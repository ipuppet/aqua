package dao

import (
	"github.com/google/wire"
)

// Set is dao providers.
var Set = wire.NewSet(
	NewRepository,
	wire.Bind(new(ArticleRepository), new(*repository)),
	wire.Bind(new(AccessRepository), new(*repository)),
	wire.Bind(new(CommentRepository), new(*repository)),
	wire.Bind(new(ManagementArticleRepository), new(*repository)),
	wire.Bind(new(ManagementTagRepository), new(*repository)),
)
