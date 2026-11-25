package logic

import (
	"context"
	"log/slog"
	"time"

	blog "aqua/services/blog/domain"
)

const searchIndexTimeout = 15 * time.Second

func (s *ManagementService) CreateArticle(ctx context.Context, article *blog.Article) (int, error) {
	aid, err := s.managementArticleRepo.CreateArticle(ctx, article)
	if err != nil {
		s.log(ctx).Error("failed to create article",
			slog.String("article.title", article.Title),
			slog.Any("error", err))
		return 0, err
	}

	s.cleanArticleCaches(ctx)
	s.log(ctx).Info("article created successfully",
		slog.Int("article.id", aid),
		slog.String("article.title", article.Title))

	article.Aid = aid
	snapshot := cloneArticle(article)
	enqueueBackground(ctx, searchIndexTimeout, "index created article", s.log, func(jobCtx context.Context) error {
		return s.indexer.IndexArticle(jobCtx, snapshot)
	})

	return aid, nil
}

func (s *ManagementService) UpdateArticle(ctx context.Context, aid int, updates *blog.ArticleUpdate) error {
	if err := s.managementArticleRepo.UpdateArticle(ctx, aid, updates); err != nil {
		s.log(ctx).Error("failed to update article",
			slog.Int("article.id", aid),
			slog.Any("error", err))
		return err
	}

	s.cleanArticleCaches(ctx)
	s.log(ctx).Info("article updated successfully",
		slog.Int("article.id", aid))

	enqueueBackground(ctx, searchIndexTimeout, "index updated article", s.log, func(jobCtx context.Context) error {
		fullArticle, err := s.articleRepo.GetArticle(jobCtx, aid)
		if err != nil {
			return err
		}
		return s.indexer.IndexArticle(jobCtx, fullArticle)
	})

	return nil
}

func (s *ManagementService) DeleteArticle(ctx context.Context, aid int) error {
	if err := s.managementArticleRepo.DeleteArticle(ctx, aid); err != nil {
		s.log(ctx).Error("failed to delete article",
			slog.Int("article.id", aid),
			slog.Any("error", err))
		return err
	}

	s.cleanArticleCaches(ctx)
	s.log(ctx).Info("article deleted successfully",
		slog.Int("article.id", aid))

	enqueueBackground(ctx, searchIndexTimeout, "delete indexed article", s.log, func(jobCtx context.Context) error {
		return s.indexer.DeleteArticle(jobCtx, aid)
	})

	return nil
}

func cloneArticle(article *blog.Article) *blog.Article {
	if article == nil {
		return nil
	}
	copied := *article
	if article.Tags != nil {
		copied.Tags = make([]blog.Tag, len(article.Tags))
		copy(copied.Tags, article.Tags)
	}
	return &copied
}
