package logic

import (
	"context"
	"log/slog"

	"github.com/ipuppet/aqua/server/pkg/infra/logutil"
	"github.com/ipuppet/aqua/server/services/blog/dao"
	"github.com/ipuppet/aqua/server/services/blog/meili"
)

// ManagementService is the logic implementation.
type ManagementService struct {
	managementArticleRepo dao.ManagementArticleRepository
	managementTagRepo     dao.ManagementTagRepository
	articleRepo           dao.ArticleRepository
	indexer               meili.BlogSearchIndexer
}

// NewManagementService composes the admin service implementation.
func NewManagementService(
	managementArticleRepo dao.ManagementArticleRepository,
	managementTagRepo dao.ManagementTagRepository,
	articleRepo dao.ArticleRepository,
	indexer meili.BlogSearchIndexer,
) *ManagementService {
	return &ManagementService{
		managementArticleRepo: managementArticleRepo,
		managementTagRepo:     managementTagRepo,
		articleRepo:           articleRepo,
		indexer:               indexer,
	}
}

func (s *ManagementService) log(ctx context.Context) *slog.Logger {
	return logutil.FromContext(ctx)
}

func (s *ManagementService) cleanArticleCaches(ctx context.Context) {
	cache := s.managementArticleRepo.Cache()
	if cache == nil {
		return
	}
	_ = cache.DelMatch(ctx, "blog.GetArticles.*")
	_ = cache.DelMatch(ctx, "blog.GetArticle.*")
	_ = cache.DelMatch(ctx, "blog.GetArticleTags.*")
	_ = cache.DelMatch(ctx, "blog.GetContext.*")
	_ = cache.DelMatch(ctx, "blog.GetCount*")
	_ = cache.DelMatch(ctx, "blog.GetLastAid*")
}

func (s *ManagementService) cleanTagCaches(ctx context.Context) {
	cache := s.managementTagRepo.Cache()
	if cache == nil {
		return
	}
	_ = cache.DelMatch(ctx, "blog.GetTags*")
	_ = cache.DelMatch(ctx, "blog.GetArticle.*")
	_ = cache.DelMatch(ctx, "blog.GetArticles.*")
	_ = cache.DelMatch(ctx, "blog.GetArticleTags.*")
}
