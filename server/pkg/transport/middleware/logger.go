package middleware

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"

	"aqua/pkg/infra/logutil"
)

// LoggerMiddleware attaches a request scoped logger enriched with trace information.
// HTTP request details are handled by otelgin.Middleware for distributed tracing.
// This middleware only injects trace context into the logger for business logic use.
func LoggerMiddleware(base *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract TraceID and SpanID from context (injected by otelgin middleware)
		spanContext := trace.SpanFromContext(c.Request.Context()).SpanContext()
		var traceID, spanID string
		if spanContext.IsValid() {
			traceID = spanContext.TraceID().String()
			spanID = spanContext.SpanID().String()
		} else {
			traceID = uuid.NewString() // Fallback if no trace context
		}

		// Create request-scoped logger with trace information.
		// HTTP request details (method, path, status, duration) are handled by otelgin.Middleware.
		reqLogger := logutil.WithTrace(base, traceID, spanID)

		// Inject logger into context for use in handlers and business logic.
		ctx := logutil.WithLogger(c.Request.Context(), reqLogger)
		c.Request = c.Request.WithContext(ctx)

		c.Set("trace_id", traceID)

		c.Next()
	}
}
