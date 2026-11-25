package imagestore

import (
	"context"
	"io"
)

// Backend defines the interface for image store backends (local disk or OpenList).
type Backend interface {
	UploadImage(ctx context.Context, module, originalName string, content io.Reader) (*ImageInfo, error)
	DeleteImage(ctx context.Context, path string) error
	ListImages(ctx context.Context, module string) ([]ImageInfo, error)
}
