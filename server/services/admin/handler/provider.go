package handler

import (
	"github.com/google/wire"
)

// Set is handler providers.
var Set = wire.NewSet(
	NewCacheHandler,
	NewFileHandler,
	wire.Struct(new(HandlerSet), "*"),
)
