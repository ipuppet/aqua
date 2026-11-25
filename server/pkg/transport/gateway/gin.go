package gateway

import (
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

// GinHandler wraps the grpc-gateway [runtime.ServeMux] and converts
// the result into a [gin.HandlerFunc].
func GinHandler(mux *runtime.ServeMux) gin.HandlerFunc {
	return gin.WrapH(mux)
}
