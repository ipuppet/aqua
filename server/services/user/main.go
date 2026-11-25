package user

import (
	"context"
	"crypto/rsa"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	userv1 "github.com/ipuppet/aqua/server/gen/go/aqua/user/v1"
	"github.com/ipuppet/aqua/server/pkg/infra/bootstrap"
	"github.com/ipuppet/aqua/server/pkg/rbac"
	"github.com/ipuppet/aqua/server/pkg/transport/gateway"
	"github.com/ipuppet/aqua/server/services/user/grpcserver"
	"github.com/ipuppet/aqua/server/services/user/logic"
)

const ServiceName = "user"

func Assemble(rt *bootstrap.Runtime) (*bootstrap.Module, error) {
	return buildModule(rt)
}

func provideModule(
	userManagementService *logic.UserManagementService,
	accessControlService *logic.AccessControlService,
	userService *logic.UserService,
	publicKey *rsa.PublicKey,
	rbacReader rbac.Reader,
) (*bootstrap.Module, error) {
	adminServer := grpcserver.NewAdminServer(userManagementService, accessControlService)
	userServer := grpcserver.NewUserServer(userService)

	gatewayMux := gateway.NewServeMux()
	if err := userv1.RegisterUserServiceHandlerServer(context.Background(), gatewayMux, userServer); err != nil {
		return nil, err
	}
	if err := userv1.RegisterAdminUserServiceHandlerServer(context.Background(), gatewayMux, adminServer); err != nil {
		return nil, err
	}

	gatewayHandler := gateway.GinHandler(gatewayMux)

	return &bootstrap.Module{
		MountHTTP: func(engine *gin.Engine) {
			registerGatewayRoutes(engine, publicKey, rbacReader, gatewayHandler)
		},
		MountGRPC: func(server *grpc.Server) {
			userv1.RegisterAdminUserServiceServer(server, adminServer)
			userv1.RegisterUserServiceServer(server, userServer)
		},
	}, nil
}
