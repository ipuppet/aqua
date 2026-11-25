package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"aqua/pkg/infra/cache"
)

var errCacheEmpty = errors.New("cache empty")

// FetchWithCache retrieves data using the Cache-Aside pattern.
// It tries to get data from the cache first. If missing, it calls fetchFunc,
// caches the result, and returns it.
func FetchWithCache(ctx context.Context, cache cache.Store, key string, ttl time.Duration, target any, fetchFunc func() (any, error)) error {
	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Ptr {
		return fmt.Errorf("target must be a pointer, got %v", targetValue.Kind())
	}
	targetElem := targetValue.Elem()

	// Try to load from cache
	if cache != nil && key != "" {
		if err := tryLoadFromCache(ctx, cache, key, target); err == nil {
			return nil
		}
		// Cache miss or error, proceed to fetch
	}

	// Fetch from database
	data, err := fetchFunc()
	if err != nil {
		if err == sql.ErrNoRows {
			return err
		}
		return fmt.Errorf("failed to fetch data: %w", err)
	}

	// Assign data to target
	if err := assignData(targetElem, data); err != nil {
		return fmt.Errorf("failed to assign data: %w", err)
	}

	// Update cache
	if cache != nil && key != "" {
		// TODO: surface cache update error via metrics or structured logging if needed.
		// For now we intentionally ignore cache write failures to keep the read path fast and resilient.
		_ = updateCache(ctx, cache, key, data, ttl)
	}

	return nil
}

func tryLoadFromCache(ctx context.Context, c cache.Store, key string, target any) error {
	cacheData, err := c.Get(ctx, key)
	if err != nil {
		return err
	}

	if len(cacheData) == 0 {
		return errCacheEmpty
	}

	if err := json.Unmarshal(cacheData, target); err != nil {
		return fmt.Errorf("failed to unmarshal cache data: %w", err)
	}

	return nil
}

func assignData(targetElem reflect.Value, data any) error {
	dataValue := reflect.ValueOf(data)

	// Direct assignment (e.g., both are slices)
	if dataValue.Type().AssignableTo(targetElem.Type()) {
		targetElem.Set(dataValue)
		return nil
	}

	// Value extraction (e.g., fetchFunc returned a pointer to struct, but targetElem is the struct value itself)
	if dataValue.Kind() == reflect.Ptr && dataValue.Elem().Type().AssignableTo(targetElem.Type()) {
		targetElem.Set(dataValue.Elem())
		return nil
	}

	return errors.New("data type mismatch")
}

func updateCache(ctx context.Context, c cache.Store, key string, data any, ttl time.Duration) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal cache data: %w", err)
	}
	// Note: We use Set with []byte directly to avoid string allocations
	return c.Set(ctx, key, jsonData, ttl)
}
