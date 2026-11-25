package grpcserver

import (
	"context"

	blogv1 "aqua/gen/go/aqua/blog/v1"
	"aqua/pkg/apperr"
	"aqua/pkg/transport/authctx"
	blog "aqua/services/blog/domain"
	"aqua/services/blog/logic"
)

type CommentServer struct {
	commentService *logic.CommentService
	blogv1.UnimplementedCommentServiceServer
}

func NewCommentServer(commentService *logic.CommentService) *CommentServer {
	return &CommentServer{commentService: commentService}
}

func (s *CommentServer) ListComments(ctx context.Context, req *blogv1.ListCommentsRequest) (*blogv1.ListCommentsResponse, error) {
	page, err := s.commentService.GetComments(ctx, int(req.Aid), int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, err
	}
	return &blogv1.ListCommentsResponse{Comments: ToPBComments(page.Comments), Total: int32(page.Total)}, nil
}

func (s *CommentServer) ListChildComments(ctx context.Context, req *blogv1.ListChildCommentsRequest) (*blogv1.ListCommentsResponse, error) {
	page, err := s.commentService.GetChildComments(ctx, int(req.Id), int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, err
	}
	return &blogv1.ListCommentsResponse{Comments: ToPBComments(page.Comments), Total: int32(page.Total)}, nil
}

func (s *CommentServer) AddComment(ctx context.Context, req *blogv1.AddCommentRequest) (*blogv1.AddCommentResponse, error) {
	tokenUser := authctx.TokenUserFromContext(ctx)
	if tokenUser == nil {
		return nil, apperr.Unauthorized("")
	}
	comment := &blog.Comment{
		Parent:  int(req.ParentId),
		Content: req.Content,
		IP:      authctx.ClientIP(ctx),
	}
	err := s.commentService.AddComment(ctx, int(req.Aid), tokenUser.UID, comment)
	return &blogv1.AddCommentResponse{}, err
}
