package middleware

import (
	"log/slog"
	"runtime/debug"

	"github.com/gin-gonic/gin"

	"aqua/pkg/apperr"
	"aqua/pkg/infra/logutil"
	"aqua/pkg/transport/response"
)

// ErrorHandler recovers panics and writes a standard JSON error envelope.
func ErrorHandler() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, err any) {
		reqLogger := logutil.FromContext(c.Request.Context())
		if reqLogger == nil {
			reqLogger = slog.Default()
		}

		reqLogger.Error("panic recovered",
			slog.Any("error", err),
			slog.String("path", c.Request.URL.Path),
			slog.String("stack", string(debug.Stack())),
		)

		response.Fail(c, apperr.Internal(""))
	})
}
