package dao

import (
	"context"
	"testing"

	"aqua/pkg/infra/cache"
)

func TestCacheRepository_FlushPreservesRefreshTokens(t *testing.T) {
	ctx := context.Background()
	m := cache.NewMemoryCache()
	repo := NewCacheRepository(m)

	// Set some keys in cache, including session keys.
	if err := m.Set(ctx, "foo", []byte("bar"), 0); err != nil {
		t.Fatalf("set foo failed: %v", err)
	}
	if err := m.Set(ctx, "refresh_token:1:abc", []byte("session"), 0); err != nil {
		t.Fatalf("set refresh_token failed: %v", err)
	}
	if err := m.Set(ctx, "baz", []byte("qux"), 0); err != nil {
		t.Fatalf("set baz failed: %v", err)
	}

	if err := repo.Flush(ctx); err != nil {
		t.Fatalf("flush failed: %v", err)
	}

	val, err := m.Get(ctx, "foo")
	if err != nil {
		t.Fatalf("get foo failed: %v", err)
	}
	if val != nil {
		t.Fatal("expected foo to be deleted")
	}

	val, err = m.Get(ctx, "baz")
	if err != nil {
		t.Fatalf("get baz failed: %v", err)
	}
	if val != nil {
		t.Fatal("expected baz to be deleted")
	}

	val, err = m.Get(ctx, "refresh_token:1:abc")
	if err != nil {
		t.Fatalf("get refresh_token failed: %v", err)
	}
	if val == nil {
		t.Fatal("expected refresh_token to be preserved")
	}
}
