package admin

import (
	"github.com/gin-gonic/gin"

	"github.com/ipuppet/aqua/server/pkg/infra/bootstrap"
	"github.com/ipuppet/aqua/server/services/admin/handler"
)

const ServiceName = "admin"

func Assemble(rt *bootstrap.Runtime) (*bootstrap.Module, error) {
	return buildModule(rt)
}

func provideModule(handlerSet handler.HandlerSet) (*bootstrap.Module, error) {
	return &bootstrap.Module{
		MountHTTP: func(engine *gin.Engine) {
			handler.RegisterRoutes(engine, handlerSet)
		},
	}, nil
}
