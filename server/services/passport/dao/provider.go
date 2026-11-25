package dao

import (
	"github.com/google/wire"
)

// Set is dao providers.
var Set = wire.NewSet(
	NewRepository,
	wire.Bind(new(AuthRepository), new(*repository)),
)
