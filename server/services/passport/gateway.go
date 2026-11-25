package passport

import (
	"crypto/rsa"

	"github.com/gin-gonic/gin"

	"github.com/ipuppet/aqua/server/pkg/rbac"
	"github.com/ipuppet/aqua/server/pkg/transport/middleware"
)

func registerGatewayRoutes(engine *gin.Engine, publicKey *rsa.PublicKey, rbacReader rbac.Reader, gatewayHandler gin.HandlerFunc) {
	if engine == nil || gatewayHandler == nil {
		return
	}

	engine.Any("/api/login/password", gatewayHandler)
	engine.Any("/api/token", gatewayHandler)
	engine.Any("/api/register", gatewayHandler)
	engine.Any("/api/verify/email", gatewayHandler)
	engine.Any("/api/password/forgot", gatewayHandler)
	engine.Any("/api/password/reset", gatewayHandler)

	sessions := engine.Group("/api/sessions")
	sessions.Use(middleware.Authenticate(publicKey))
	sessions.Any("", gatewayHandler)
	sessions.Any("/*any", gatewayHandler)

	adminSessions := engine.Group("/api/users")
	adminSessions.Use(middleware.Authenticate(publicKey))
	adminSessions.Use(middleware.RequireAnyPermission(rbacReader, "user.manager"))
	adminSessions.Any("/:uid/sessions", gatewayHandler)
	adminSessions.Any("/:uid/sessions/*any", gatewayHandler)
}
