package logic

import (
	"context"
	"log/slog"

	"github.com/ipuppet/aqua/server/pkg/infra/logutil"
	"github.com/ipuppet/aqua/server/services/blog/dao"
	blog "github.com/ipuppet/aqua/server/services/blog/domain"
)

// CommentService is the logic implementation.
type CommentService struct {
	repo dao.CommentRepository
}

// NewCommentService builds a comment service backed by the repository implementation.
func NewCommentService(repo dao.CommentRepository) *CommentService {
	return &CommentService{repo: repo}
}

func (s *CommentService) log(ctx context.Context) *slog.Logger {
	return logutil.FromContext(ctx)
}

func (s *CommentService) GetComments(ctx context.Context, aid, page, pageSize int) (*blog.CommentPage, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	comments, err := s.repo.GetComments(ctx, aid, offset, pageSize)
	if err != nil {
		s.log(ctx).Error("failed to get comments", "aid", aid, "error", err)
		return nil, err
	}

	total, err := s.repo.GetCommentCount(ctx, aid)
	if err != nil {
		s.log(ctx).Error("failed to get comment count", "aid", aid, "error", err)
		return nil, err
	}

	return &blog.CommentPage{Comments: comments, Total: total}, nil
}

func (s *CommentService) GetChildComments(ctx context.Context, id, page, pageSize int) (*blog.CommentPage, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	comments, err := s.repo.GetChildComments(ctx, id, offset, pageSize)
	if err != nil {
		s.log(ctx).Error("failed to get child comments", "parent", id, "error", err)
		return nil, err
	}

	total, err := s.repo.GetChildCommentCount(ctx, id)
	if err != nil {
		s.log(ctx).Error("failed to get child comment count", "parent", id, "error", err)
		return nil, err
	}

	return &blog.CommentPage{Comments: comments, Total: total}, nil
}

func (s *CommentService) AddComment(ctx context.Context, aid, uid int, comment *blog.Comment) error {
	comment.UID = uid
	comment.Aid = aid
	if err := s.repo.AddComment(ctx, comment); err != nil {
		s.log(ctx).Error("failed to add comment", "aid", aid, "error", err)
		return err
	}

	s.log(ctx).Info("added comment", "aid", aid, "uid", comment.UID)
	return nil
}
