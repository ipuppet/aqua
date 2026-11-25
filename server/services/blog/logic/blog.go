package logic

import (
	"context"
	"log/slog"
	"time"

	"github.com/ipuppet/aqua/server/pkg/infra/logutil"
	"github.com/ipuppet/aqua/server/services/blog/dao"
	blog "github.com/ipuppet/aqua/server/services/blog/domain"
	"github.com/ipuppet/aqua/server/services/blog/meili"
)

const accessLogTimeout = 5 * time.Second

// ArticleService is the default ArticleService implementation.
type ArticleService struct {
	repo       dao.ArticleRepository
	accessRepo dao.AccessRepository
	indexer    meili.BlogSearchIndexer
}

// NewArticleService composes the ArticleService implementation.
func NewArticleService(repo dao.ArticleRepository, accessRepo dao.AccessRepository, indexer meili.BlogSearchIndexer) *ArticleService {
	return &ArticleService{repo: repo, accessRepo: accessRepo, indexer: indexer}
}

func (s *ArticleService) log(ctx context.Context) *slog.Logger {
	return logutil.FromContext(ctx)
}

func (s *ArticleService) GetArticle(ctx context.Context, aid int) (*blog.Article, error) {
	article, err := s.repo.GetArticle(ctx, aid)
	if err != nil {
		s.log(ctx).Error("failed to get article",
			slog.Int("article.id", aid),
			slog.Any("error", err))
		return nil, err
	}
	return article, nil
}

func (s *ArticleService) GetArticleDetail(ctx context.Context, aid int, ip string) (*blog.ArticleDetail, error) {
	logger := s.log(ctx).With(slog.Int("article.id", aid))
	lastAid, err := s.repo.GetLastAid(ctx)
	if err != nil {
		logger.Error("failed to get last article ID",
			slog.Any("error", err))
		return nil, err
	}

	previousArticle, err := s.repo.GetContext(ctx, aid, -1, lastAid)
	if err != nil {
		logger.Warn("failed to get previous article context",
			slog.Any("error", err))
	}

	nextArticle, err := s.repo.GetContext(ctx, aid, 1, lastAid)
	if err != nil {
		logger.Warn("failed to get next article context",
			slog.Any("error", err))
	}

	article, err := s.repo.GetArticle(ctx, aid)
	if err != nil {
		logger.Error("failed to get article",
			slog.Any("error", err))
		return nil, err
	}

	enqueueBackground(ctx, accessLogTimeout, "log article access", s.log, func(jobCtx context.Context) error {
		return s.accessRepo.LogAccessInfo(jobCtx, aid, ip)
	})

	return &blog.ArticleDetail{
		Article:  article,
		Previous: previousArticle,
		Next:     nextArticle,
	}, nil
}

func (s *ArticleService) GetArticles(ctx context.Context, page, pageSize int) (*blog.ArticlesResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	articles, err := s.repo.GetArticles(ctx, offset, pageSize)
	if err != nil {
		s.log(ctx).Error("failed to get articles",
			slog.Int("page", page),
			slog.Int("page_size", pageSize),
			slog.Any("error", err))
		return nil, err
	}

	total, err := s.repo.GetCount(ctx)
	if err != nil {
		s.log(ctx).Error("failed to get article count",
			slog.Any("error", err))
		return nil, err
	}

	return &blog.ArticlesResponse{
		Articles: articles,
		Total:    total,
	}, nil
}

func (s *ArticleService) SearchArticles(ctx context.Context, kw string, page, pageSize int) (*blog.SearchResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	searchResp, err := s.indexer.SearchArticles(ctx, kw, pageSize, offset)
	if err != nil {
		s.log(ctx).Error("failed to search articles via meilisearch",
			slog.String("search.keyword", kw),
			slog.Any("error", err))
		return nil, err
	}

	return searchResp, nil
}

func (s *ArticleService) GetTags(ctx context.Context) ([]blog.Tag, error) {
	tags, err := s.repo.GetTags(ctx)
	if err != nil {
		s.log(ctx).Error("failed to get article tags",
			slog.Any("error", err))
		return nil, err
	}
	return tags, nil
}

func (s *ArticleService) GetAbout(ctx context.Context) (string, error) {
	article, err := s.repo.GetArticle(ctx, -1)
	if err != nil {
		s.log(ctx).Error("failed to get about page content",
			slog.Any("error", err))
		return "", err
	}
	return article.Markdown, nil
}

func (s *ArticleService) GetAccess(ctx context.Context, day int) ([]blog.AccessResponse, error) {
	accesses, err := s.accessRepo.GetAccess(ctx, day)
	if err != nil {
		s.log(ctx).Error("failed to get article access statistics",
			slog.Int("stats.days", day),
			slog.Any("error", err))
		return nil, err
	}
	responses := make([]blog.AccessResponse, 0, len(accesses))
	for _, access := range accesses {
		responses = append(responses, *ToAccessResponse(&access))
	}
	return responses, nil
}

func (s *ArticleService) GetArticlesByTagID(ctx context.Context, id int) ([]blog.ArticleDescription, error) {
	articles, err := s.repo.GetArticlesByTagID(ctx, id)
	if err != nil {
		s.log(ctx).Error("failed to get articles by tag",
			slog.Int("tag.id", id),
			slog.Any("error", err))
		return nil, err
	}
	return articles, nil
}
