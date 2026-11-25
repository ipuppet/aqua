package database

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"

	"aqua/pkg/infra/redis"
)

func TestFetchWithCache(t *testing.T) {
	db, mock := redismock.NewClientMock()
	cache := redis.NewRedisCache(db)
	ctx := context.Background()

	// Case 1: Cache hit
	mock.ExpectGet("hit").SetVal("\"cached_value\"")
	var target string
	err := FetchWithCache(ctx, cache, "hit", time.Minute, &target, func() (any, error) {
		return "db_value", nil
	})
	assert.NoError(t, err)
	assert.Equal(t, "cached_value", target)

	// Case 2: Cache miss, DB success
	mock.ExpectGet("miss").SetErr(errors.New("redis: nil"))
	mock.ExpectSet("miss", []byte("\"db_value\""), time.Minute).SetVal("OK")
	err = FetchWithCache(ctx, cache, "miss", time.Minute, &target, func() (any, error) {
		return "db_value", nil
	})
	assert.NoError(t, err)
	assert.Equal(t, "db_value", target)

	// Case 3: Cache miss, DB error
	mock.ExpectGet("error").SetErr(errors.New("redis: nil"))
	err = FetchWithCache(ctx, cache, "error", time.Minute, &target, func() (any, error) {
		return nil, errors.New("db error")
	})
	assert.Error(t, err)

	// Case 4: Cache miss, DB success, cache set fails (ignored)
	mock.ExpectGet("set_fail").SetErr(errors.New("redis: nil"))
	mock.ExpectSet("set_fail", []byte("\"db_value\""), time.Minute).SetErr(errors.New("set failed"))
	err = FetchWithCache(ctx, cache, "set_fail", time.Minute, &target, func() (any, error) {
		return "db_value", nil
	})
	assert.NoError(t, err)
	assert.Equal(t, "db_value", target)
}
