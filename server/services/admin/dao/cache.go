package dao

import (
	"context"
	"strings"

	"github.com/ipuppet/aqua/server/pkg/infra/cache"
)

type cacheRepository struct {
	cache cache.Store
}

var _ CacheRepository = (*cacheRepository)(nil)

// NewCacheRepository creates cache repository implementation.
func NewCacheRepository(c cache.Store) CacheRepository {
	return &cacheRepository{cache: c}
}

func (r *cacheRepository) Flush(ctx context.Context) error {
	if r.cache == nil {
		return nil
	}

	// Keep user sessions (`refresh_token:*`) during admin cache clear.
	allKeys, err := r.cache.GetMatch(ctx, "*")
	if err != nil {
		return err
	}

	for key := range allKeys {
		if strings.HasPrefix(key, "refresh_token:") {
			continue
		}
		if err := r.cache.Del(ctx, key); err != nil {
			return err
		}
	}

	return nil
}
