package meili

import (
	"github.com/google/wire"
)

// Set provides BlogSearchIndexer.
var Set = wire.NewSet(NewBlogSearchIndexer)
