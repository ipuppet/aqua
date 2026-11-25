package __SERVICE__

import (
	"github.com/gin-gonic/gin"
)

func registerGatewayRoutes(engine *gin.Engine, gatewayHandler gin.HandlerFunc) {
	if engine == nil || gatewayHandler == nil {
		return
	}

	engine.Any("/api/__SERVICE__/hello", gatewayHandler)
}
