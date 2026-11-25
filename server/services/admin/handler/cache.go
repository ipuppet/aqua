package handler

import (
	"crypto/rsa"

	"github.com/gin-gonic/gin"

	"aqua/pkg/transport/middleware"
	"aqua/pkg/transport/response"
	"aqua/services/admin/dao"
)

// CacheHandler exposes cache management endpoints.
type CacheHandler struct {
	cacheService dao.CacheRepository
	publicKey    *rsa.PublicKey
}

// NewCacheHandler constructs a CacheHandler.
func NewCacheHandler(cacheService dao.CacheRepository, publicKey *rsa.PublicKey) *CacheHandler {
	return &CacheHandler{
		cacheService: cacheService,
		publicKey:    publicKey,
	}
}

// Register wires cache routes.
func (h *CacheHandler) Register(e *gin.Engine) {
	if h == nil {
		return
	}
	r := e.Group("/api")
	r.Use(middleware.Authenticate(h.publicKey))
	r.DELETE("/cache", func(gc *gin.Context) {
		err := h.cacheService.Flush(gc.Request.Context())
		response.OKOrFail(gc, err)
	})
}
