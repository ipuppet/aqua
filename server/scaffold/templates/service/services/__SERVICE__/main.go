package __SERVICE__

import (
	"context"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	__SERVICE__v1 "aqua/gen/go/aqua/__SERVICE__/v1"
	"aqua/pkg/infra/bootstrap"
	"aqua/pkg/transport/gateway"
	"aqua/services/__SERVICE__/grpcserver"
	"aqua/services/__SERVICE__/logic"
)

const ServiceName = "__SERVICE__"

func Assemble(rt *bootstrap.Runtime) (*bootstrap.Module, error) {
	return buildModule(rt)
}

func provideModule(
	helloService *logic.HelloService,
) (*bootstrap.Module, error) {
	server := grpcserver.New__Service__Server(helloService)

	gatewayMux := gateway.NewServeMux()
	if err := __SERVICE__v1.Register__Service__ServiceHandlerServer(context.Background(), gatewayMux, server); err != nil {
		return nil, err
	}

	gatewayHandler := gateway.GinHandler(gatewayMux)

	return &bootstrap.Module{
		MountHTTP: func(engine *gin.Engine) {
			registerGatewayRoutes(engine, gatewayHandler)
		},
		MountGRPC: func(gs *grpc.Server) {
			__SERVICE__v1.Register__Service__ServiceServer(gs, server)
		},
	}, nil
}
