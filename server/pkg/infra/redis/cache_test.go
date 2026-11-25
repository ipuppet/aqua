package redis

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

func TestRedisCache_Get(t *testing.T) {
	db, mock := redismock.NewClientMock()
	cache := NewRedisCache(db)
	ctx := context.Background()

	mock.ExpectGet("key").SetVal("value")

	val, err := cache.Get(ctx, "key")
	assert.NoError(t, err)
	assert.Equal(t, []byte("value"), val)

	mock.ExpectGet("missing").SetErr(errors.New("redis error"))
	_, err = cache.Get(ctx, "missing")
	assert.Error(t, err)
}

func TestRedisCache_Set(t *testing.T) {
	db, mock := redismock.NewClientMock()
	cache := NewRedisCache(db)
	ctx := context.Background()

	mock.ExpectSet("key", "value", time.Minute).SetVal("OK")

	err := cache.Set(ctx, "key", "value", time.Minute)
	assert.NoError(t, err)

	mock.ExpectSet("error", "value", time.Minute).SetErr(errors.New("redis error"))
	err = cache.Set(ctx, "error", "value", time.Minute)
	assert.Error(t, err)
}

func TestRedisCache_Del(t *testing.T) {
	db, mock := redismock.NewClientMock()
	cache := NewRedisCache(db)
	ctx := context.Background()

	mock.ExpectDel("key").SetVal(1)

	err := cache.Del(ctx, "key")
	assert.NoError(t, err)

	mock.ExpectDel("error").SetErr(errors.New("redis error"))
	err = cache.Del(ctx, "error")
	assert.Error(t, err)
}

func TestRedisCache_DelMatch(t *testing.T) {
	db, mock := redismock.NewClientMock()
	cache := NewRedisCache(db)
	ctx := context.Background()

	// Mock Scan
	mock.ExpectScan(0, "prefix*", 0).SetVal([]string{"key1", "key2"}, 0)
	mock.ExpectDel("key1").SetVal(1)
	mock.ExpectDel("key2").SetVal(1)

	err := cache.DelMatch(ctx, "prefix*")
	assert.NoError(t, err)

	// Mock Scan error
	mock.ExpectScan(0, "error*", 0).SetErr(errors.New("scan error"))
	err = cache.DelMatch(ctx, "error*")
	assert.Error(t, err)
}
