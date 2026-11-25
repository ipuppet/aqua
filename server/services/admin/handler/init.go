package handler

import "github.com/gin-gonic/gin"

// HandlerSet aggregates admin HTTP handlers that are not served via grpc-gateway.
type HandlerSet struct {
	Cache *CacheHandler
	File  *FileHandler
}

// RegisterRoutes registers file and cache handlers.
func RegisterRoutes(e *gin.Engine, set HandlerSet) {
	if set.Cache != nil {
		set.Cache.Register(e)
	}
	if set.File != nil {
		set.File.Register(e)
	}
}
