package user

import (
	"crypto/rsa"

	"github.com/gin-gonic/gin"

	"aqua/pkg/rbac"
	"aqua/pkg/transport/middleware"
)

func registerGatewayRoutes(engine *gin.Engine, publicKey *rsa.PublicKey, rbacReader rbac.Reader, gatewayHandler gin.HandlerFunc) {
	if engine == nil || gatewayHandler == nil {
		return
	}

	me := engine.Group("/api/users/me")
	me.Use(middleware.Authenticate(publicKey))
	me.Any("", gatewayHandler)
	me.Any("/*any", gatewayHandler)

	users := engine.Group("/api/users")
	users.Use(middleware.Authenticate(publicKey))
	users.Use(middleware.RequireAnyPermission(rbacReader, "user.manager"))
	users.Any("", gatewayHandler)
	users.Any("/search/*any", gatewayHandler)
	users.PUT("/:uid", gatewayHandler)
	users.Any("/:uid/roles", gatewayHandler)
	users.Any("/:uid/roles/*any", gatewayHandler)
	users.Any("/:uid/permissions", gatewayHandler)

	engine.GET("/api/users/:uid", gatewayHandler)

	access := engine.Group("/api/access-control")
	access.Use(middleware.Authenticate(publicKey))
	access.Use(middleware.RequireAnyPermission(rbacReader, "rbac.manager"))
	access.Any("", gatewayHandler)
	access.Any("/*any", gatewayHandler)
}
