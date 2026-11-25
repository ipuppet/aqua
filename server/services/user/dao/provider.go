package dao

import (
	"github.com/google/wire"
)

// Set is dao providers.
var Set = wire.NewSet(
	NewRepository,
	wire.Bind(new(UserRepository), new(*repository)),
	wire.Bind(new(UserManagementRepository), new(*repository)),
	wire.Bind(new(AccessControlRepository), new(*repository)),
)
