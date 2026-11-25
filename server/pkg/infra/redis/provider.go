package redis

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/wire"
	"github.com/redis/go-redis/v9"

	"aqua/pkg/infra/cache"
	"aqua/pkg/infra/config"
	"aqua/pkg/infra/models"
)

// redisPingTimeout 是每次探活 Ping 的最长等待时间，同时被 reconnect.go 中 tryConnect 使用。
const redisPingTimeout = 3 * time.Second

// NewCache returns a cache.Store backed by Redis.
func NewCache(watcher *config.Watcher[*models.AppConfig], logger *slog.Logger) (cache.Store, error) {
	if watcher == nil {
		return nil, fmt.Errorf("redis watcher is not configured")
	}
	appCfg := watcher.Value()
	if appCfg == nil {
		return nil, fmt.Errorf("app config is not loaded")
	}

	client, err := newRedisClient(appCfg.Redis)
	if err != nil {
		return nil, err
	}

	redisCache := NewRedisCache(client)
	watcher.OnChange(func(nextAppCfg *models.AppConfig) {
		if nextAppCfg == nil {
			return
		}
		nextClient, err := newRedisClient(nextAppCfg.Redis)
		if err != nil {
			if logger != nil {
				logger.Warn("redis configuration reload failed", slog.Any("error", err))
			}
			return
		}
		redisCache.UpdateClient(nextClient)
		if logger != nil {
			logger.Info("redis configuration reloaded successfully")
		}
	})

	return redisCache, nil
}

func newRedisClient(redisCfg models.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Addr,
		Password: redisCfg.Password,
		DB:       redisCfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), redisPingTimeout)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("connect redis: %w", err)
	}

	return client, nil
}

// Set provides cache.Store.
var Set = wire.NewSet(NewCache)
