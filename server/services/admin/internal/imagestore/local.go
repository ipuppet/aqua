package imagestore

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"aqua/pkg/infra/logutil"
	"aqua/pkg/utils/crypto"
	"aqua/services/admin/types"
)

// localBackend stores images on local disk.
type localBackend struct {
	config types.FileStoreConfig
	clock  func() time.Time
	logger func(context.Context) *slog.Logger
}

var _ Backend = (*localBackend)(nil)

// NewLocalBackend creates a local-disk-backed image store backend.
func NewLocalBackend(cfg types.FileStoreConfig, clock func() time.Time) Backend {
	return &localBackend{config: cfg, clock: clock, logger: logutil.FromContext}
}

func (b *localBackend) UploadImage(ctx context.Context, module, originalName string, content io.Reader) (*ImageInfo, error) {
	if err := validateSegment(module); err != nil {
		return nil, err
	}
	if originalName == "" {
		return nil, errors.New("original file name is required")
	}
	if content == nil {
		return nil, errors.New("file content is required")
	}

	extName := filepath.Ext(originalName)
	fileName := crypto.MD5(originalName) + extName
	dateDir := b.clock().Format("2006-01-02")
	relPath := path.Join(module, dateDir, fileName)

	saveDir := filepath.Join(b.config.SavePath, module, dateDir)
	if err := os.MkdirAll(saveDir, 0o755); err != nil {
		b.logger(ctx).Error("failed to create upload directory",
			slog.String("file.path", saveDir),
			slog.Any("error", err))
		return nil, err
	}

	dstPath := filepath.Join(saveDir, fileName)
	dst, err := os.Create(dstPath)
	if err != nil {
		b.logger(ctx).Error("failed to create file",
			slog.String("file.path", dstPath),
			slog.Any("error", err))
		return nil, err
	}
	defer func() { _ = dst.Close() }()

	written, err := io.Copy(dst, content)
	if err != nil {
		b.logger(ctx).Error("failed to save uploaded file",
			slog.String("file.path", dstPath),
			slog.Any("error", err))
		return nil, err
	}

	baseURL := strings.TrimRight(b.config.URLPrefix, "/")
	imageURL := relPath
	if baseURL != "" {
		imageURL = baseURL + "/" + strings.TrimLeft(relPath, "/")
	}

	return &ImageInfo{
		Path:    relPath,
		Name:    fileName,
		URL:     imageURL,
		Size:    written,
		ModTime: b.clock(),
	}, nil
}

func (b *localBackend) DeleteImage(ctx context.Context, path string) error {
	if path == "" {
		return errors.New("path is required")
	}
	if err := validatePath(path); err != nil {
		return err
	}

	target := filepath.Join(b.config.SavePath, path)
	if err := os.Remove(target); err != nil {
		b.logger(ctx).Error("failed to delete file",
			slog.String("file.path", target),
			slog.Any("error", err))
		return err
	}
	return nil
}

func (b *localBackend) ListImages(ctx context.Context, module string) ([]ImageInfo, error) {
	if err := validateSegment(module); err != nil {
		return nil, err
	}

	moduleDir := filepath.Join(b.config.SavePath, module)
	if _, err := os.Stat(moduleDir); err != nil {
		if os.IsNotExist(err) {
			return []ImageInfo{}, nil
		}
		b.logger(ctx).Error("failed to stat module directory",
			slog.String("file.path", moduleDir),
			slog.Any("error", err))
		return nil, err
	}

	imageExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
		".webp": true, ".svg": true, ".bmp": true, ".ico": true, ".avif": true,
	}

	var images []ImageInfo
	if err := b.collectImages(ctx, moduleDir, module, imageExts, &images); err != nil {
		return nil, err
	}

	sort.Slice(images, func(i, j int) bool {
		if images[i].ModTime != images[j].ModTime {
			return images[i].ModTime.After(images[j].ModTime)
		}
		return images[i].Path < images[j].Path
	})

	return images, nil
}

// collectImages recursively walks dirPath, collecting image files into images.
// relPrefix is the relative path from store root (e.g. module or module/subdir).
func (b *localBackend) collectImages(ctx context.Context, dirPath, relPrefix string, imageExts map[string]bool, images *[]ImageInfo) error {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		b.logger(ctx).Warn("failed to read directory",
			slog.String("file.path", dirPath),
			slog.Any("error", err))
		return err
	}

	for _, entry := range entries {
		if err := validateSegment(entry.Name()); err != nil {
			b.logger(ctx).Warn("skipping invalid entry",
				slog.String("file.name", entry.Name()),
				slog.String("file.path", dirPath),
				slog.Any("error", err))
			continue
		}

		if entry.IsDir() {
			subPath := filepath.Join(dirPath, entry.Name())
			subRel := path.Join(relPrefix, entry.Name())
			if err := b.collectImages(ctx, subPath, subRel, imageExts, images); err != nil {
				slog.Warn("failed to collect images from subdirectory", slog.String("subPath", subPath), slog.Any("error", err))
			}
			continue
		}

		fileName := entry.Name()
		ext := strings.ToLower(filepath.Ext(fileName))
		if !imageExts[ext] {
			continue
		}

		filePath := filepath.Join(dirPath, fileName)
		info, err := entry.Info()
		if err != nil {
			b.logger(ctx).Warn("failed to get file information",
				slog.String("file.path", filePath),
				slog.Any("error", err))
			continue
		}

		relPath := path.Join(relPrefix, fileName)
		baseURL := strings.TrimRight(b.config.URLPrefix, "/")
		imageURL := relPath
		if baseURL != "" {
			imageURL = baseURL + "/" + strings.TrimLeft(relPath, "/")
		}

		*images = append(*images, ImageInfo{
			Path:    relPath,
			Name:    fileName,
			URL:     imageURL,
			Size:    info.Size(),
			ModTime: info.ModTime(),
		})
	}
	return nil
}
