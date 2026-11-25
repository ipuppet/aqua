package dao

import (
	"context"

	"github.com/ipuppet/aqua/server/pkg/infra/cache"
	"github.com/ipuppet/aqua/server/pkg/infra/database"
	blog "github.com/ipuppet/aqua/server/services/blog/domain"
)

// ArticleRepository exposes read-side article operations.
type ArticleRepository interface {
	Cache() cache.Store
	GetArticle(ctx context.Context, aid int) (*blog.Article, error)
	ListAllArticleIDs(ctx context.Context) ([]int, error)
	GetLastAid(ctx context.Context) (int, error)
	GetContext(ctx context.Context, startAid, direction, boundary int) (*blog.ArticleContext, error)
	GetArticles(ctx context.Context, offset, limit int) ([]blog.Article, error)
	GetCount(ctx context.Context) (int, error)
	GetTags(ctx context.Context) ([]blog.Tag, error)
	GetArticlesByTagID(ctx context.Context, id int) ([]blog.ArticleDescription, error)
}

// AccessRepository exposes article access logging and statistics operations.
type AccessRepository interface {
	GetAccess(ctx context.Context, day int) ([]blog.Access, error)
	LogAccessInfo(ctx context.Context, aid int, ip string) error
}

// ManagementArticleRepository exposes write operations for blog articles.
type ManagementArticleRepository interface {
	Cache() cache.Store
	CreateArticle(ctx context.Context, article *blog.Article) (int, error)
	CreateArticleWithAid(ctx context.Context, article *blog.Article) error
	UpdateArticle(ctx context.Context, aid int, updates *blog.ArticleUpdate) error
	DeleteArticle(ctx context.Context, aid int) error
}

// ManagementTagRepository exposes write operations for blog tags.
type ManagementTagRepository interface {
	Cache() cache.Store
	CreateTag(ctx context.Context, name string) error
	UpdateTag(ctx context.Context, id int, updates *blog.TagUpdate) error
	DeleteTag(ctx context.Context, id int) error
}

// CommentRepository exposes comment operations.
type CommentRepository interface {
	Cache() cache.Store
	GetComments(ctx context.Context, aid, offset, limit int) ([]blog.Comment, error)
	GetCommentCount(ctx context.Context, aid int) (int, error)
	GetChildComments(ctx context.Context, id, offset, limit int) ([]blog.Comment, error)
	GetChildCommentCount(ctx context.Context, id int) (int, error)
	AddComment(ctx context.Context, comment *blog.Comment) error
}

// repository implements all blog data interfaces.
type repository struct {
	db    *database.Client
	cache cache.Store
}

// NewRepository builds a repository backed by the provided database connection.
func NewRepository(db *database.Client, c cache.Store) *repository {
	return &repository{db: db, cache: c}
}

func (r *repository) Cache() cache.Store {
	return r.cache
}

var _ ArticleRepository = (*repository)(nil)
var _ AccessRepository = (*repository)(nil)
var _ ManagementArticleRepository = (*repository)(nil)
var _ ManagementTagRepository = (*repository)(nil)
var _ CommentRepository = (*repository)(nil)
