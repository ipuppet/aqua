package rbac

import "github.com/google/wire"

// Set provides RBAC read dependencies.
var Set = wire.NewSet(NewReader)

// EventBusSet provides RBAC event invalidation dependencies.
var EventBusSet = wire.NewSet(
	NewCacheInvalidationHandler,
	NewEventBus,
)
