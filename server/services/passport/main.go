package passport

import (
	"context"
	"crypto/rsa"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	passportv1 "github.com/ipuppet/aqua/server/gen/go/aqua/passport/v1"
	"github.com/ipuppet/aqua/server/pkg/infra/bootstrap"
	"github.com/ipuppet/aqua/server/pkg/rbac"
	"github.com/ipuppet/aqua/server/pkg/transport/gateway"
	"github.com/ipuppet/aqua/server/services/passport/grpcserver"
	"github.com/ipuppet/aqua/server/services/passport/logic"
)

const ServiceName = "passport"

func Assemble(rt *bootstrap.Runtime) (*bootstrap.Module, error) {
	return buildModule(rt)
}

func provideModule(
	authService *logic.AuthService,
	registration *logic.RegistrationService,
	passwordReset *logic.PasswordResetService,
	sessionService *logic.SessionService,
	publicKey *rsa.PublicKey,
	rbacReader rbac.Reader,
) (*bootstrap.Module, error) {
	adminServer := grpcserver.NewAdminServer(sessionService)
	passportServer := grpcserver.NewPassportServer(
		authService,
		registration,
		passwordReset,
		sessionService,
	)

	gatewayMux := gateway.NewServeMux()
	if err := passportv1.RegisterPassportServiceHandlerServer(context.Background(), gatewayMux, passportServer); err != nil {
		return nil, err
	}
	if err := passportv1.RegisterAdminPassportServiceHandlerServer(context.Background(), gatewayMux, adminServer); err != nil {
		return nil, err
	}

	gatewayHandler := gateway.GinHandler(gatewayMux)

	return &bootstrap.Module{
		MountHTTP: func(engine *gin.Engine) {
			registerGatewayRoutes(engine, publicKey, rbacReader, gatewayHandler)
		},
		MountGRPC: func(server *grpc.Server) {
			passportv1.RegisterAdminPassportServiceServer(server, adminServer)
			passportv1.RegisterPassportServiceServer(server, passportServer)
		},
	}, nil
}
