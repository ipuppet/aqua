package rbac

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type cacheRecorder struct {
	delKeys     []string
	delPatterns []string
}

func (c *cacheRecorder) Get(ctx context.Context, key string) ([]byte, error) { return nil, nil }
func (c *cacheRecorder) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return nil
}
func (c *cacheRecorder) Del(ctx context.Context, key string) error {
	c.delKeys = append(c.delKeys, key)
	return nil
}
func (c *cacheRecorder) DelMatch(ctx context.Context, pattern string) error {
	c.delPatterns = append(c.delPatterns, pattern)
	return nil
}
func (c *cacheRecorder) GetMatch(ctx context.Context, pattern string) (map[string][]byte, error) {
	return make(map[string][]byte), nil
}
func (c *cacheRecorder) Flush(ctx context.Context) error { return nil }

func TestCacheInvalidationHandler_PermissionEvents(t *testing.T) {
	recorder := &cacheRecorder{}
	h := NewCacheInvalidationHandler(recorder)
	ctx := context.Background()

	h.Handle(ctx, PermissionUpdated{})
	h.Handle(ctx, PermissionDeleted{})

	assert.Equal(t, []string{RolePermissionsPattern, RolePermissionsPattern}, recorder.delPatterns)
	assert.Empty(t, recorder.delKeys)
}

func TestCacheInvalidationHandler_RolePermissionEvents(t *testing.T) {
	recorder := &cacheRecorder{}
	h := NewCacheInvalidationHandler(recorder)
	ctx := context.Background()

	h.Handle(ctx, RolePermissionAssigned{Role: "admin"})
	h.Handle(ctx, RolePermissionRevoked{Role: "owner"})

	assert.Equal(t, []string{CacheKeyStringRolePermissions("admin"), CacheKeyStringRolePermissions("owner")}, recorder.delKeys)
}

func TestCacheInvalidationHandler_RoleUpdated(t *testing.T) {
	recorder := &cacheRecorder{}
	h := NewCacheInvalidationHandler(recorder)
	ctx := context.Background()

	h.Handle(ctx, RoleUpdated{Role: "auditor"})

	assert.Equal(t, []string{CacheKeyStringRolePermissions("auditor")}, recorder.delKeys)
	assert.Equal(t, []string{UserRolesPattern}, recorder.delPatterns)
}
