// Package persist exports the shared database + redis Wire set for domain services.
package persist

import (
	"github.com/google/wire"

	"github.com/ipuppet/aqua/server/pkg/infra/database"
	"github.com/ipuppet/aqua/server/pkg/infra/redis"
)

// Set provides database and redis for persistence-backed services.
var Set = wire.NewSet(
	database.Set,
	redis.Set,
)
