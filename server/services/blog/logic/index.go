package logic

import (
	"context"
	"log/slog"
)

func (s *ManagementService) SyncMeiliIndex(ctx context.Context) (int, error) {
	if err := s.indexer.EnsureBlogIndex(); err != nil {
		s.log(ctx).Error("failed to ensure meilisearch index", slog.Any("error", err))
		return 0, err
	}

	ids, err := s.articleRepo.ListAllArticleIDs(ctx)
	if err != nil {
		s.log(ctx).Error("failed to list article ids", slog.Any("error", err))
		return 0, err
	}

	count := 0
	for _, aid := range ids {
		article, err := s.articleRepo.GetArticle(ctx, aid)
		if err != nil {
			s.log(ctx).Warn("failed to get article for indexing", slog.Int("aid", aid), slog.Any("error", err))
			continue
		}
		if err := s.indexer.IndexArticle(ctx, article); err != nil {
			s.log(ctx).Warn("failed to index article", slog.Int("aid", aid), slog.Any("error", err))
			continue
		}
		count++
	}

	s.log(ctx).Info("meilisearch sync completed", slog.Int("total_indexed", count))
	return count, nil
}
