package blog

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

	engine.GET("/api/articles", gatewayHandler)
	engine.GET("/api/articles/:aid", gatewayHandler)
	engine.GET("/api/search/*any", gatewayHandler)
	engine.GET("/api/archive", gatewayHandler)
	engine.GET("/api/tags", gatewayHandler)
	engine.GET("/api/tags/:id/articles", gatewayHandler)
	engine.GET("/api/about", gatewayHandler)

	engine.GET("/api/comments/articles/:aid", gatewayHandler)
	engine.GET("/api/comments/:id/children", gatewayHandler)

	commentWrite := engine.Group("/api/comments/articles")
	commentWrite.Use(middleware.Authenticate(publicKey))
	commentWrite.POST("/:aid", gatewayHandler)

	articleWrite := engine.Group("/api/articles")
	articleWrite.Use(middleware.Authenticate(publicKey))
	articleWrite.Use(middleware.RequireAnyPermission(rbacReader, "blog.manager"))
	articleWrite.POST("", gatewayHandler)
	articleWrite.GET("/:aid/edit", gatewayHandler)
	articleWrite.PUT("/:aid", gatewayHandler)
	articleWrite.DELETE("/:aid", gatewayHandler)

	tagWrite := engine.Group("/api/tags")
	tagWrite.Use(middleware.Authenticate(publicKey))
	tagWrite.Use(middleware.RequireAnyPermission(rbacReader, "blog.manager"))
	tagWrite.POST("", gatewayHandler)
	tagWrite.PUT("/:id", gatewayHandler)
	tagWrite.DELETE("/:id", gatewayHandler)

	aboutWrite := engine.Group("/api/about")
	aboutWrite.Use(middleware.Authenticate(publicKey))
	aboutWrite.Use(middleware.RequireAnyPermission(rbacReader, "blog.manager"))
	aboutWrite.PUT("", gatewayHandler)

	admin := engine.Group("/api")
	admin.Use(middleware.Authenticate(publicKey))
	admin.Use(middleware.RequireAnyPermission(rbacReader, "blog.manager"))
	admin.GET("/access/:day", gatewayHandler)
	admin.POST("/meili/sync", gatewayHandler)
}
