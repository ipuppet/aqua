package blog

import (
	"context"
	"crypto/rsa"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	blogv1 "github.com/ipuppet/aqua/server/gen/go/aqua/blog/v1"
	"github.com/ipuppet/aqua/server/pkg/infra/bootstrap"
	"github.com/ipuppet/aqua/server/pkg/rbac"
	"github.com/ipuppet/aqua/server/pkg/transport/gateway"
	"github.com/ipuppet/aqua/server/services/blog/grpcserver"
	"github.com/ipuppet/aqua/server/services/blog/logic"
	blogmeili "github.com/ipuppet/aqua/server/services/blog/meili"
)

const ServiceName = "blog"

func Assemble(rt *bootstrap.Runtime) (*bootstrap.Module, error) {
	return buildModule(rt)
}

func provideModule(
	rt *bootstrap.Runtime,
	articleService *logic.ArticleService,
	commentService *logic.CommentService,
	blogManagementService *logic.ManagementService,
	blogSearchIndexer blogmeili.BlogSearchIndexer,
	publicKey *rsa.PublicKey,
	rbacReader rbac.Reader,
) (*bootstrap.Module, error) {
	if err := blogSearchIndexer.EnsureBlogIndex(); err != nil {
		rt.Logger.Error("failed to ensure meilisearch index", "error", err)
	}

	adminServer := grpcserver.NewAdminServer(articleService, blogManagementService)
	articleServer := grpcserver.NewArticleServer(articleService)
	commentServer := grpcserver.NewCommentServer(commentService)

	gatewayMux := gateway.NewServeMux()
	if err := blogv1.RegisterArticleServiceHandlerServer(context.Background(), gatewayMux, articleServer); err != nil {
		return nil, err
	}
	if err := blogv1.RegisterCommentServiceHandlerServer(context.Background(), gatewayMux, commentServer); err != nil {
		return nil, err
	}
	if err := blogv1.RegisterAdminBlogServiceHandlerServer(context.Background(), gatewayMux, adminServer); err != nil {
		return nil, err
	}

	gatewayHandler := gateway.GinHandler(gatewayMux)

	return &bootstrap.Module{
		MountHTTP: func(engine *gin.Engine) {
			registerGatewayRoutes(engine, publicKey, rbacReader, gatewayHandler)
		},
		MountGRPC: func(server *grpc.Server) {
			blogv1.RegisterAdminBlogServiceServer(server, adminServer)
			blogv1.RegisterArticleServiceServer(server, articleServer)
			blogv1.RegisterCommentServiceServer(server, commentServer)
		},
	}, nil
}
