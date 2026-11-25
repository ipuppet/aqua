package rbac

import (
	"context"
	"sync"

	"aqua/pkg/infra/cache"
)

// Event is a marker interface for RBAC mutations that require cache invalidation.
type Event interface {
	isRBACEvent()
}

type UserRoleAssigned struct {
	UID int
}

func (UserRoleAssigned) isRBACEvent() {}

type UserRoleRevoked struct {
	UID int
}

func (UserRoleRevoked) isRBACEvent() {}

type RolePermissionAssigned struct {
	Role string
}

func (RolePermissionAssigned) isRBACEvent() {}

type RolePermissionRevoked struct {
	Role string
}

func (RolePermissionRevoked) isRBACEvent() {}

type RoleUpdated struct {
	Role string
}

func (RoleUpdated) isRBACEvent() {}

type RoleDeleted struct {
	Role string
}

func (RoleDeleted) isRBACEvent() {}

type PermissionUpdated struct{}

func (PermissionUpdated) isRBACEvent() {}

type PermissionDeleted struct{}

func (PermissionDeleted) isRBACEvent() {}

// EventBus dispatches RBAC events to subscribers.
type EventBus interface {
	Publish(ctx context.Context, event Event)
	Subscribe(handler Handler)
}

// Handler consumes RBAC events.
type Handler interface {
	Handle(ctx context.Context, event Event)
}

type eventBus struct {
	handlers []Handler
	mu       sync.RWMutex
}

// NewEventBus creates an in-process event bus.
func NewEventBus(handler Handler) EventBus {
	handlers := []Handler{}
	if handler != nil {
		handlers = append(handlers, handler)
	}
	return &eventBus{handlers: handlers}
}

func (b *eventBus) Publish(ctx context.Context, event Event) {
	if event == nil {
		return
	}

	b.mu.RLock()
	handlers := append([]Handler(nil), b.handlers...)
	b.mu.RUnlock()

	for _, h := range handlers {
		h.Handle(ctx, event)
	}
}

func (b *eventBus) Subscribe(handler Handler) {
	if handler == nil {
		return
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers = append(b.handlers, handler)
}

// CacheInvalidationHandler invalidates RBAC cache keys based on events.
type CacheInvalidationHandler struct {
	cache cache.Store
}

func NewCacheInvalidationHandler(c cache.Store) Handler {
	return &CacheInvalidationHandler{cache: c}
}

func (h *CacheInvalidationHandler) Handle(ctx context.Context, event Event) {
	if h.cache == nil {
		return
	}
	switch e := event.(type) {
	case UserRoleAssigned:
		_ = h.cache.Del(ctx, CacheKeyUserRoles(e.UID))
	case UserRoleRevoked:
		_ = h.cache.Del(ctx, CacheKeyUserRoles(e.UID))
	case RolePermissionAssigned:
		_ = h.cache.Del(ctx, CacheKeyStringRolePermissions(e.Role))
	case RolePermissionRevoked:
		_ = h.cache.Del(ctx, CacheKeyStringRolePermissions(e.Role))
	case RoleUpdated:
		_ = h.cache.DelMatch(ctx, UserRolesPattern)
		_ = h.cache.Del(ctx, CacheKeyStringRolePermissions(e.Role))
	case RoleDeleted:
		_ = h.cache.DelMatch(ctx, UserRolesPattern)
		_ = h.cache.Del(ctx, CacheKeyStringRolePermissions(e.Role))
	case PermissionUpdated:
		_ = h.cache.DelMatch(ctx, RolePermissionsPattern)
	case PermissionDeleted:
		_ = h.cache.DelMatch(ctx, RolePermissionsPattern)
	}
}
