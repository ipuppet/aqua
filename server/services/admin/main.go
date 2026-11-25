package admin

import (
	"github.com/gin-gonic/gin"

	"aqua/pkg/infra/bootstrap"
	"aqua/services/admin/handler"
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
