package logutil

import (
	"context"
	"log/slog"
)

// ctxKey represents the context key for storing request-scoped loggers.
type ctxKey struct{}

// WithLogger returns a new context containing the provided logger.
// If ctx is nil, it defaults to context.Background().
// If logger is nil, it returns the original context unmodified.
func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if logger == nil {
		return ctx
	}
	return context.WithValue(ctx, ctxKey{}, logger)
}

// FromContext returns the logger stored in the context.
// If no logger is found or ctx is nil, it returns the default logger.
func FromContext(ctx context.Context) *slog.Logger {
	if ctx != nil {
		if logger, ok := ctx.Value(ctxKey{}).(*slog.Logger); ok && logger != nil {
			return logger
		}
	}
	return Default()
}
