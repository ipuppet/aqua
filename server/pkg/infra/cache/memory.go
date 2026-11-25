package cache

import (
	"context"
	"encoding/json"
	"path/filepath"
	"sync"
	"time"
)

// entry holds cached value and expiration time.
type entry struct {
	data     []byte
	expireAt time.Time
}

// MemoryCache implements Store using an in-memory map.
// Used as fallback when Redis is unavailable.
// DelMatch uses filepath.Match for glob patterns (e.g. "prefix*").
var _ Store = (*MemoryCache)(nil)

type MemoryCache struct {
	data map[string]entry
	mu   sync.RWMutex
}

// NewMemoryCache creates a new in-memory cache.
func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		data: make(map[string]entry),
	}
}

// Get retrieves a value. Returns nil, nil on cache miss.
// Expired entries are removed lazily and treated as miss.
func (m *MemoryCache) Get(ctx context.Context, key string) ([]byte, error) {
	m.mu.RLock()
	e, ok := m.data[key]
	m.mu.RUnlock()

	if !ok {
		return nil, nil
	}
	if !e.expireAt.IsZero() && time.Now().After(e.expireAt) {
		_ = m.Del(ctx, key)
		return nil, nil
	}
	return e.data, nil
}

// Set stores a value. Accepts []byte directly or any value (JSON-marshaled).
func (m *MemoryCache) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	var data []byte
	switch v := value.(type) {
	case []byte:
		data = v
	default:
		var err error
		data, err = json.Marshal(value)
		if err != nil {
			return err
		}
	}

	expireAt := time.Time{}
	if expiration > 0 {
		expireAt = time.Now().Add(expiration)
	}

	m.mu.Lock()
	m.data[key] = entry{data: data, expireAt: expireAt}
	m.mu.Unlock()
	return nil
}

// Del removes a key.
func (m *MemoryCache) Del(ctx context.Context, key string) error {
	m.mu.Lock()
	delete(m.data, key)
	m.mu.Unlock()
	return nil
}

// DelMatch removes keys matching the glob pattern (e.g. "prefix*").
func (m *MemoryCache) DelMatch(ctx context.Context, pattern string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for key := range m.data {
		ok, _ := filepath.Match(pattern, key)
		if ok {
			delete(m.data, key)
		}
	}
	return nil
}

// GetMatch returns all key-value pairs matching the glob pattern.
func (m *MemoryCache) GetMatch(ctx context.Context, pattern string) (map[string][]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make(map[string][]byte)
	for key, e := range m.data {
		ok, _ := filepath.Match(pattern, key)
		if !ok {
			continue
		}
		// Skip expired entries
		if !e.expireAt.IsZero() && time.Now().After(e.expireAt) {
			continue
		}
		result[key] = e.data
	}
	return result, nil
}

// Flush removes all keys.
func (m *MemoryCache) Flush(ctx context.Context) error {
	m.mu.Lock()
	m.data = make(map[string]entry)
	m.mu.Unlock()
	return nil
}
