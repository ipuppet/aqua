package grpcserver

import (
	"context"

	blogv1 "github.com/ipuppet/aqua/server/gen/go/aqua/blog/v1"
	"github.com/ipuppet/aqua/server/pkg/apperr"
	"github.com/ipuppet/aqua/server/pkg/transport/authctx"
	"github.com/ipuppet/aqua/server/services/blog/logic"
)

type ArticleServer struct {
	articleService *logic.ArticleService
	blogv1.UnimplementedArticleServiceServer
}

func NewArticleServer(articleService *logic.ArticleService) *ArticleServer {
	return &ArticleServer{articleService: articleService}
}

func (s *ArticleServer) ListArticles(ctx context.Context, req *blogv1.ListArticlesRequest) (*blogv1.ListArticlesResponse, error) {
	data, err := s.articleService.GetArticles(ctx, int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, err
	}
	return &blogv1.ListArticlesResponse{Articles: ToPBArticles(data.Articles), Total: int32(data.Total)}, nil
}

func (s *ArticleServer) GetArticleDetail(ctx context.Context, req *blogv1.GetArticleRequest) (*blogv1.GetArticleDetailResponse, error) {
	detail, err := s.articleService.GetArticleDetail(ctx, int(req.Aid), authctx.ClientIP(ctx))
	if err != nil {
		return nil, apperr.FromDAO(err)
	}

	return &blogv1.GetArticleDetailResponse{
		Article:  ToPBArticle(detail.Article),
		Previous: ToPBArticleContext(detail.Previous),
		Next:     ToPBArticleContext(detail.Next),
	}, nil
}

func (s *ArticleServer) SearchArticles(ctx context.Context, req *blogv1.SearchArticlesRequest) (*blogv1.SearchArticlesResponse, error) {
	data, err := s.articleService.SearchArticles(ctx, req.Keyword, int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, err
	}
	return &blogv1.SearchArticlesResponse{Articles: ToPBArticles(data.Items), Total: int32(data.Total)}, nil
}

func (s *ArticleServer) GetArchive(ctx context.Context, req *blogv1.ListArticlesRequest) (*blogv1.ListArticlesResponse, error) {
	data, err := s.articleService.GetArticles(ctx, 1, 1000)
	if err != nil {
		return nil, err
	}
	return &blogv1.ListArticlesResponse{Articles: ToPBArticles(data.Articles), Total: int32(data.Total)}, nil
}

func (s *ArticleServer) ListTags(ctx context.Context, req *blogv1.ListTagsRequest) (*blogv1.ListTagsResponse, error) {
	tags, err := s.articleService.GetTags(ctx)
	if err != nil {
		return nil, err
	}
	return &blogv1.ListTagsResponse{Tags: ToPBTags(tags)}, nil
}

func (s *ArticleServer) ListArticlesByTag(ctx context.Context, req *blogv1.ListArticlesByTagRequest) (*blogv1.ListArticlesByTagResponse, error) {
	data, err := s.articleService.GetArticlesByTagID(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}
	return &blogv1.ListArticlesByTagResponse{Articles: ToPBArticleDescriptions(data)}, nil
}

func (s *ArticleServer) GetAbout(ctx context.Context, req *blogv1.GetAboutRequest) (*blogv1.GetAboutResponse, error) {
	content, err := s.articleService.GetAbout(ctx)
	if err != nil {
		return nil, apperr.FromDAO(err)
	}
	return &blogv1.GetAboutResponse{Content: content}, nil
}
