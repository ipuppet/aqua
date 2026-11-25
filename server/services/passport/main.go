package passport

import (
	"context"
	"crypto/rsa"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	passportv1 "aqua/gen/go/aqua/passport/v1"
	"aqua/pkg/infra/bootstrap"
	"aqua/pkg/rbac"
	"aqua/pkg/transport/gateway"
	"aqua/services/passport/grpcserver"
	"aqua/services/passport/logic"
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
