package dao

import "github.com/google/wire"

// Set is dao providers.
var Set = wire.NewSet(
	NewCacheRepository,
)
