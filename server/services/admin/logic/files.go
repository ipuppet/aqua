package logic

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"time"

	"aqua/pkg/infra/cache"
	"aqua/services/admin/internal/imagestore"
)

const (
	fileListImagesCachePrefix = "file:ListImages:"
	fileListImagesCacheTTL    = time.Hour
)

// FileService handles uploads and deletions within the configured store root.
type FileService struct {
	configProvider FileConfigProvider
	store          cache.Store
}

// NewFileService constructs a FileService.
func NewFileService(configProvider FileConfigProvider, store cache.Store) *FileService {
	return &FileService{configProvider: configProvider, store: store}
}

func (s *FileService) cache() cache.Store {
	return s.store
}

func (s *FileService) getBackendAndCache() (imagestore.Backend, bool, error) {
	cfg, err := s.configProvider.FileStoreConfig()
	if err != nil {
		return nil, false, err
	}
	var backend imagestore.Backend
	useCache := false
	if cfg.OpenList != nil && cfg.OpenList.Enable && cfg.OpenList.Endpoint != "" {
		backend = imagestore.NewOpenlistBackend(*cfg.OpenList, time.Now)
		useCache = true
	} else {
		backend = imagestore.NewLocalBackend(cfg, time.Now)
	}
	return backend, useCache, nil
}

func (s *FileService) retrieveFromCache(ctx context.Context, cacheKey string) ([]imagestore.ImageInfo, error) {
	if c := s.cache(); c != nil {
		raw, err := c.Get(ctx, cacheKey)
		if err != nil {
			return nil, err
		}
		var images []imagestore.ImageInfo
		if err := json.Unmarshal([]byte(raw), &images); err != nil {
			return nil, err
		}
		return images, nil
	}
	return nil, errors.New("cache disabled")
}

func (s *FileService) storeToCache(ctx context.Context, cacheKey string, images []imagestore.ImageInfo) error {
	if c := s.cache(); c != nil {
		if data, err := json.Marshal(images); err == nil {
			_ = c.Set(ctx, cacheKey, string(data), fileListImagesCacheTTL)
		}
	}
	return nil
}

func (s *FileService) updateCache(ctx context.Context, cacheKey string, update func([]imagestore.ImageInfo) []imagestore.ImageInfo, retrive func() ([]imagestore.ImageInfo, error)) {
	if s.cache() == nil {
		return
	}
	// 尝试乐观更新
	images, err := s.retrieveFromCache(ctx, cacheKey)
	if err == nil {
		// 命中缓存才进行乐观更新，避免在 Cache Miss 的情况下将单张图写为全量缓存
		images = update(images)
		_ = s.storeToCache(ctx, cacheKey, images)
	}

	// 全量回源保证数据一致性（最终一致）
	newImages, err := retrive()
	if err != nil {
		return
	}
	_ = s.storeToCache(ctx, cacheKey, newImages)
}

func (s *FileService) ListImages(ctx context.Context, module string, refresh bool) ([]imagestore.ImageInfo, error) {
	backend, useCache, err := s.getBackendAndCache()
	if err != nil {
		return nil, err
	}

	cacheKey := fileListImagesCachePrefix + module

	if useCache {
		if refresh {
			if c := s.cache(); c != nil {
				_ = c.Del(ctx, cacheKey)
			}
		} else {
			images, err := s.retrieveFromCache(ctx, cacheKey)
			if err == nil {
				return images, nil
			}
		}
	}

	images, err := backend.ListImages(ctx, module)
	if err != nil {
		return nil, err
	}

	if useCache {
		_ = s.storeToCache(ctx, cacheKey, images)
	}
	return images, nil
}

func (s *FileService) UploadImage(ctx context.Context, module, originalName string, content io.Reader) (*imagestore.ImageInfo, error) {
	backend, useCache, err := s.getBackendAndCache()
	if err != nil {
		return nil, err
	}

	info, err := backend.UploadImage(ctx, module, originalName, content)
	if err != nil {
		return nil, err
	}
	if useCache {
		cacheKey := fileListImagesCachePrefix + module
		bgCtx := context.WithoutCancel(ctx)
		go s.updateCache(
			bgCtx,
			cacheKey,
			func(images []imagestore.ImageInfo) []imagestore.ImageInfo {
				images = append(images, *info)
				return images
			},
			func() ([]imagestore.ImageInfo, error) {
				return backend.ListImages(bgCtx, module)
			},
		)
	}
	return info, nil
}

func (s *FileService) DeleteImage(ctx context.Context, path string) error {
	backend, useCache, err := s.getBackendAndCache()
	if err != nil {
		return err
	}

	err = backend.DeleteImage(ctx, path)
	if err != nil {
		return err
	}

	if useCache {
		module := path[:strings.Index(path, "/")]
		cacheKey := fileListImagesCachePrefix + module
		bgCtx := context.WithoutCancel(ctx)
		go s.updateCache(
			bgCtx,
			cacheKey,
			func(images []imagestore.ImageInfo) []imagestore.ImageInfo {
				for i, image := range images {
					if image.Path == path {
						images = append(images[:i], images[i+1:]...)
						break
					}
				}
				return images
			},
			func() ([]imagestore.ImageInfo, error) {
				return backend.ListImages(bgCtx, module)
			},
		)
	}

	return nil
}
