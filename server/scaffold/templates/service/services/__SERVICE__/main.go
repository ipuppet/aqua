package __SERVICE__

import (
	"context"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	__SERVICE__v1 "github.com/ipuppet/aqua/server/gen/go/aqua/__SERVICE__/v1"
	"github.com/ipuppet/aqua/server/pkg/infra/bootstrap"
	"github.com/ipuppet/aqua/server/pkg/transport/gateway"
	"github.com/ipuppet/aqua/server/services/__SERVICE__/grpcserver"
	"github.com/ipuppet/aqua/server/services/__SERVICE__/logic"
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
