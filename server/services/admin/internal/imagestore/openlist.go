package imagestore

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/ipuppet/aqua/server/pkg/infra/logutil"
	"github.com/ipuppet/aqua/server/pkg/utils/crypto"
	"github.com/ipuppet/aqua/server/services/admin/pkg/openlist"
	"github.com/ipuppet/aqua/server/services/admin/types"
)

// openlistBackend stores images on OpenList.
type openlistBackend struct {
	client *openlist.Client
	cfg    types.OpenListConfig
	clock  func() time.Time
	logger func(context.Context) *slog.Logger
}

var _ Backend = (*openlistBackend)(nil)

// NewOpenlistBackend creates an OpenList-backed image store backend.
func NewOpenlistBackend(cfg types.OpenListConfig, clock func() time.Time) Backend {
	return &openlistBackend{
		client: openlist.NewClient(cfg),
		cfg:    cfg,
		clock:  clock,
		logger: logutil.FromContext,
	}
}

func (b *openlistBackend) UploadImage(ctx context.Context, module, originalName string, content io.Reader) (*ImageInfo, error) {
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

	destPath := "/" + path.Join(strings.Trim(b.cfg.PathPrefix, "/"), module, dateDir, fileName)
	imageURL, err := b.client.Upload(ctx, destPath, content)
	if err != nil {
		b.logger(ctx).Error("failed to upload to openlist", slog.Any("error", err))
		return nil, err
	}
	return &ImageInfo{
		Path:    relPath,
		Name:    fileName,
		URL:     imageURL,
		Size:    0,
		ModTime: b.clock(),
	}, nil
}

func (b *openlistBackend) DeleteImage(ctx context.Context, relPath string) error {
	if relPath == "" {
		return errors.New("path is required")
	}
	if err := validatePath(relPath); err != nil {
		return err
	}

	dirPath := path.Dir(relPath)
	fileName := path.Base(relPath)
	baseDir := path.Join(strings.Trim(b.cfg.PathPrefix, "/"), dirPath)
	dir := "/" + baseDir
	return b.client.Remove(ctx, dir, []string{fileName})
}

func (b *openlistBackend) ListImages(ctx context.Context, module string) ([]ImageInfo, error) {
	if err := validateSegment(module); err != nil {
		return nil, err
	}

	imageExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
		".webp": true, ".svg": true, ".bmp": true, ".ico": true, ".avif": true,
	}
	base := path.Join(strings.Trim(b.cfg.PathPrefix, "/"), module)
	dirPath := "/" + base

	var images []ImageInfo
	if err := b.collectImages(ctx, dirPath, path.Join(module), imageExts, &images); err != nil {
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
// relPrefix is the relative path from module root (e.g. module or module/subdir).
func (b *openlistBackend) collectImages(ctx context.Context, dirPath, relPrefix string, imageExts map[string]bool, images *[]ImageInfo) error {
	const perPage = 100
	var content []openlist.FsObject
	for page := 1; ; page++ {
		chunk, total, err := b.client.List(ctx, dirPath, page, perPage)
		if err != nil {
			b.logger(ctx).Error("failed to list directory",
				slog.String("file.path", dirPath),
				slog.Any("error", err))
			return err
		}
		content = append(content, chunk...)
		if len(content) >= total || len(chunk) < perPage {
			break
		}
	}

	for _, entry := range content {
		if err := validateSegment(entry.Name); err != nil {
			b.logger(ctx).Warn("skipping invalid entry",
				slog.String("file.name", entry.Name),
				slog.String("file.path", dirPath),
				slog.Any("error", err))
			continue
		}

		if entry.IsDir {
			subPath := path.Join(dirPath, entry.Name)
			subRel := path.Join(relPrefix, entry.Name)
			if err := b.collectImages(ctx, subPath, subRel, imageExts, images); err != nil {
				b.logger(ctx).Warn("failed to list subdirectory",
					slog.String("file.path", subPath),
					slog.Any("error", err))
				// continue with other entries instead of returning
			}
			continue
		}

		ext := strings.ToLower(filepath.Ext(entry.Name))
		if !imageExts[ext] {
			continue
		}

		relPath := path.Join(relPrefix, entry.Name)
		modTime, _ := time.Parse(time.RFC3339, entry.Modified)
		imageURL := buildDirectURL(b.cfg.Endpoint, b.cfg.PathPrefix, relPath)

		*images = append(*images, ImageInfo{
			Path:    relPath,
			Name:    entry.Name,
			URL:     imageURL,
			Size:    entry.Size,
			ModTime: modTime,
		})
	}
	return nil
}
