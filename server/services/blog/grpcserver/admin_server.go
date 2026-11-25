package grpcserver

import (
	"context"

	blogv1 "github.com/ipuppet/aqua/server/gen/go/aqua/blog/v1"
	"github.com/ipuppet/aqua/server/pkg/apperr"
	blog "github.com/ipuppet/aqua/server/services/blog/domain"
	"github.com/ipuppet/aqua/server/services/blog/logic"
)

// AdminServer provides gRPC APIs for admin-side blog operations.
type AdminServer struct {
	articleService    *logic.ArticleService
	managementService *logic.ManagementService
	blogv1.UnimplementedAdminBlogServiceServer
}

// NewAdminServer constructs AdminServer.
func NewAdminServer(
	articleService *logic.ArticleService,
	managementService *logic.ManagementService,
) *AdminServer {
	return &AdminServer{
		articleService:    articleService,
		managementService: managementService,
	}
}

func (s *AdminServer) GetArticle(ctx context.Context, req *blogv1.GetArticleRequest) (*blogv1.GetArticleResponse, error) {
	item, err := s.articleService.GetArticle(ctx, int(req.GetAid()))
	if err != nil {
		return nil, apperr.FromDAO(err)
	}
	return &blogv1.GetArticleResponse{Article: ToPBArticle(item)}, nil
}

func (s *AdminServer) CreateArticle(ctx context.Context, req *blogv1.CreateArticleRequest) (*blogv1.CreateArticleResponse, error) {
	aid, err := s.managementService.CreateArticle(ctx, &blog.Article{
		Title:    req.GetTitle(),
		Summary:  req.GetSummary(),
		Markdown: req.GetMarkdown(),
		Tags:     ToModelTags(req.GetTags()),
	})
	if err != nil {
		return nil, err
	}
	return &blogv1.CreateArticleResponse{Aid: int32(aid)}, nil
}

func (s *AdminServer) UpdateArticle(ctx context.Context, req *blogv1.UpdateArticleRequest) (*blogv1.UpdateArticleResponse, error) {
	updates := &blog.ArticleUpdate{}
	if req.Title != nil {
		v := req.Title.GetValue()
		updates.Title = &v
	}
	if req.Summary != nil {
		v := req.Summary.GetValue()
		updates.Summary = &v
	}
	if req.Markdown != nil {
		v := req.Markdown.GetValue()
		updates.Markdown = &v
	}
	// HTTP admin always sends the full article, including tags.
	tags := ToModelTags(req.GetTags())
	updates.Tags = &tags
	if err := s.managementService.UpdateArticle(ctx, int(req.GetAid()), updates); err != nil {
		return nil, err
	}
	return &blogv1.UpdateArticleResponse{}, nil
}

func (s *AdminServer) DeleteArticle(ctx context.Context, req *blogv1.DeleteArticleRequest) (*blogv1.DeleteArticleResponse, error) {
	if err := s.managementService.DeleteArticle(ctx, int(req.GetAid())); err != nil {
		return nil, err
	}
	return &blogv1.DeleteArticleResponse{}, nil
}

func (s *AdminServer) CreateTag(ctx context.Context, req *blogv1.CreateTagRequest) (*blogv1.CreateTagResponse, error) {
	if err := s.managementService.CreateTag(ctx, req.GetName()); err != nil {
		return nil, err
	}
	return &blogv1.CreateTagResponse{}, nil
}

func (s *AdminServer) UpdateTag(ctx context.Context, req *blogv1.UpdateTagRequest) (*blogv1.UpdateTagResponse, error) {
	name := req.GetName()
	if err := s.managementService.UpdateTag(ctx, int(req.GetId()), &blog.TagUpdate{Name: &name}); err != nil {
		return nil, err
	}
	return &blogv1.UpdateTagResponse{}, nil
}

func (s *AdminServer) DeleteTag(ctx context.Context, req *blogv1.DeleteTagRequest) (*blogv1.DeleteTagResponse, error) {
	if err := s.managementService.DeleteTag(ctx, int(req.GetId())); err != nil {
		return nil, err
	}
	return &blogv1.DeleteTagResponse{}, nil
}

func (s *AdminServer) SetAbout(ctx context.Context, req *blogv1.SetAboutRequest) (*blogv1.SetAboutResponse, error) {
	if err := s.managementService.UpdateAbout(ctx, req.GetContent()); err != nil {
		return nil, err
	}
	return &blogv1.SetAboutResponse{}, nil
}

func (s *AdminServer) GetAccessStats(ctx context.Context, req *blogv1.GetAccessStatsRequest) (*blogv1.GetAccessStatsResponse, error) {
	items, err := s.articleService.GetAccess(ctx, int(req.GetDay()))
	if err != nil {
		return nil, apperr.FromDAO(err)
	}
	return &blogv1.GetAccessStatsResponse{Records: ToPBAccess(items)}, nil
}

func (s *AdminServer) SyncMeili(ctx context.Context, _ *blogv1.SyncMeiliRequest) (*blogv1.SyncMeiliResponse, error) {
	count, err := s.managementService.SyncMeiliIndex(ctx)
	if err != nil {
		return nil, err
	}
	return &blogv1.SyncMeiliResponse{TotalIndexed: int32(count)}, nil
}
