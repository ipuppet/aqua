// Package cache provides a generic cache interface and implementations.
// Redis 不可用时可降级为 MemoryCache。
package cache

import (
	"context"
	"time"
)

// Store defines the interface for caching operations.
type Store interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	Del(ctx context.Context, key string) error
	DelMatch(ctx context.Context, pattern string) error
	GetMatch(ctx context.Context, pattern string) (map[string][]byte, error)
	Flush(ctx context.Context) error
}
