package discovery

import "github.com/google/wire"

// Set is discovery providers.
var Set = wire.NewSet(NewEtcdServiceRegistry)
