package database

import (
	"github.com/google/wire"
)

// Set provides a reloadable database client.
var Set = wire.NewSet(NewClient)
