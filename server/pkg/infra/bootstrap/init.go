// Package bootstrap provides application initialization, lifecycle management, and server startup routines.
package bootstrap

import "github.com/gin-gonic/gin"

func init() {
	gin.SetMode(gin.ReleaseMode)
}
