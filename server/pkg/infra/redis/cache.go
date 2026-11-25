package redis

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"

	"aqua/pkg/infra/cache"
)

// RedisCache implements cache.Store using Redis
type RedisCache struct {
	mu     sync.RWMutex
	client *redis.Client
}

var _ cache.Store = (*RedisCache)(nil)

// NewRedisCache creates a RedisCache with the given client.
// Intended for testing or scenarios requiring a custom redis.Client.
func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{client: client}
}

// UpdateClient swaps the underlying Redis client and closes the previous one.
func (c *RedisCache) UpdateClient(client *redis.Client) {
	if client == nil {
		return
	}

	c.mu.Lock()
	old := c.client
	c.client = client
	c.mu.Unlock()

	if old != nil {
		_ = old.Close()
	}
}

func (c *RedisCache) currentClient() (*redis.Client, error) {
	c.mu.RLock()
	client := c.client
	c.mu.RUnlock()
	if client == nil {
		return nil, fmt.Errorf("redis client is not initialized")
	}
	return client, nil
}

// Get gets a value from Redis
func (c *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	client, err := c.currentClient()
	if err != nil {
		return nil, err
	}
	return client.Get(ctx, key).Bytes()
}

// Set sets a value in Redis
func (c *RedisCache) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	client, err := c.currentClient()
	if err != nil {
		return err
	}
	return client.Set(ctx, key, value, expiration).Err()
}

// Del deletes a value from Redis
func (c *RedisCache) Del(ctx context.Context, key string) error {
	client, err := c.currentClient()
	if err != nil {
		return err
	}
	return client.Del(ctx, key).Err()
}

// DelMatch deletes keys matching the pattern
func (c *RedisCache) DelMatch(ctx context.Context, pattern string) error {
	client, err := c.currentClient()
	if err != nil {
		return err
	}

	iter := client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := client.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}

// GetMatch returns all key-value pairs matching the pattern
func (c *RedisCache) GetMatch(ctx context.Context, pattern string) (map[string][]byte, error) {
	client, err := c.currentClient()
	if err != nil {
		return nil, err
	}

	result := make(map[string][]byte)
	iter := client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		val, err := client.Get(ctx, key).Bytes()
		if err != nil && err != redis.Nil {
			return nil, err
		}
		if err == nil {
			result[key] = val
		}
	}
	return result, iter.Err()
}

// Flush flushes all keys from the Redis database
func (c *RedisCache) Flush(ctx context.Context) error {
	client, err := c.currentClient()
	if err != nil {
		return err
	}
	return client.FlushDB(ctx).Err()
}
