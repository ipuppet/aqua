package logic

import (
	"context"
	"log/slog"

	blog "github.com/ipuppet/aqua/server/services/blog/domain"
)

func (s *ManagementService) CreateTag(ctx context.Context, name string) error {
	if err := s.managementTagRepo.CreateTag(ctx, name); err != nil {
		s.log(ctx).Error("failed to create tag",
			slog.String("tag.name", name),
			slog.Any("error", err))
		return err
	}

	s.cleanTagCaches(ctx)
	s.log(ctx).Info("tag created successfully",
		slog.String("tag.name", name))
	return nil
}

func (s *ManagementService) UpdateTag(ctx context.Context, id int, updates *blog.TagUpdate) error {
	if err := s.managementTagRepo.UpdateTag(ctx, id, updates); err != nil {
		s.log(ctx).Error("failed to update tag",
			slog.Int("tag.id", id),
			slog.Any("error", err))
		return err
	}

	s.cleanTagCaches(ctx)
	s.log(ctx).Info("tag updated successfully",
		slog.Int("tag.id", id))
	return nil
}

func (s *ManagementService) DeleteTag(ctx context.Context, id int) error {
	if err := s.managementTagRepo.DeleteTag(ctx, id); err != nil {
		s.log(ctx).Error("failed to delete tag",
			slog.Int("tag.id", id),
			slog.Any("error", err))
		return err
	}

	s.cleanTagCaches(ctx)
	s.log(ctx).Info("tag deleted successfully",
		slog.Int("tag.id", id))
	return nil
}
