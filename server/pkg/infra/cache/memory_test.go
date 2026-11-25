package cache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMemoryCache_GetSetDel(t *testing.T) {
	ctx := context.Background()
	m := NewMemoryCache()

	err := m.Set(ctx, "k1", []byte("v1"), time.Minute)
	assert.NoError(t, err)

	val, err := m.Get(ctx, "k1")
	assert.NoError(t, err)
	assert.Equal(t, []byte("v1"), val)

	err = m.Del(ctx, "k1")
	assert.NoError(t, err)

	val, err = m.Get(ctx, "k1")
	assert.NoError(t, err)
	assert.Nil(t, val)
}

func TestMemoryCache_DelMatch(t *testing.T) {
	ctx := context.Background()
	m := NewMemoryCache()

	_ = m.Set(ctx, "blog:GetArticle:1", []byte("a"), time.Minute)
	_ = m.Set(ctx, "blog:GetArticle:2", []byte("b"), time.Minute)
	_ = m.Set(ctx, "other:key", []byte("c"), time.Minute)

	err := m.DelMatch(ctx, "blog:GetArticle:*")
	assert.NoError(t, err)

	val, _ := m.Get(ctx, "blog:GetArticle:1")
	assert.Nil(t, val)
	val, _ = m.Get(ctx, "blog:GetArticle:2")
	assert.Nil(t, val)
	val, _ = m.Get(ctx, "other:key")
	assert.Equal(t, []byte("c"), val)
}

func TestMemoryCache_Flush(t *testing.T) {
	ctx := context.Background()
	m := NewMemoryCache()

	_ = m.Set(ctx, "k1", []byte("v1"), time.Minute)
	err := m.Flush(ctx)
	assert.NoError(t, err)

	val, _ := m.Get(ctx, "k1")
	assert.Nil(t, val)
}
