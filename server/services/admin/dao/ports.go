package dao

import (
	"context"
)

// CacheRepository defines cache management operations.
type CacheRepository interface {
	Flush(ctx context.Context) error
}
