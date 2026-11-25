// Package persist exports the shared database + redis Wire set for domain services.
package persist

import (
	"github.com/google/wire"

	"aqua/pkg/infra/database"
	"aqua/pkg/infra/redis"
)

// Set provides database and redis for persistence-backed services.
var Set = wire.NewSet(
	database.Set,
	redis.Set,
)
