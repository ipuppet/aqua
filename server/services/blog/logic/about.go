package logic

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	blog "aqua/services/blog/domain"
)

func (s *ManagementService) UpdateAbout(ctx context.Context, content string) error {
	if err := s.ensureAboutArticle(ctx, content); err != nil {
		return err
	}

	title := "About"
	summary := "About me..."
	emptyTags := []blog.Tag{}
	updates := &blog.ArticleUpdate{
		Title:    &title,
		Markdown: &content,
		Summary:  &summary,
		Tags:     &emptyTags,
	}

	if err := s.managementArticleRepo.UpdateArticle(ctx, -1, updates); err != nil {
		s.log(ctx).Error("failed to update about article",
			slog.Any("error", err))
		return err
	}

	if cache := s.managementArticleRepo.Cache(); cache != nil {
		_ = cache.DelMatch(ctx, "blog.GetArticle.-1*")
	}
	s.log(ctx).Info("about article updated successfully")
	return nil
}

func (s *ManagementService) ensureAboutArticle(ctx context.Context, content string) error {
	_, err := s.articleRepo.GetArticle(ctx, -1)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		t := time.Date(2017, 10, 1, 0, 0, 0, 0, time.UTC)
		article := &blog.Article{
			Aid:       -1,
			Title:     "About",
			Markdown:  content,
			Summary:   "About me...",
			Tags:      []blog.Tag{},
			CreatedAt: &t,
		}
		if err := s.managementArticleRepo.CreateArticleWithAid(ctx, article); err != nil {
			s.log(ctx).Error("failed to create about article",
				slog.Any("error", err))
			return err
		}
		s.log(ctx).Info("about article created successfully")
		return nil
	case err != nil:
		s.log(ctx).Error("failed to check about article existence",
			slog.Any("error", err))
		return err
	default:
		return nil
	}
}
