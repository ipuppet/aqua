package logutil

import (
	"context"
	"log/slog"
)

// MultiHandler represents a handler that duplicates log records to multiple underlying handlers.
// It is safe for concurrent use provided the configured handlers are also safely concurrent.
type MultiHandler struct {
	handlers []slog.Handler
}

// NewMultiHandler returns an initialized MultiHandler wrapping the provided handlers.
// If no handlers are provided, it returns an empty MultiHandler.
func NewMultiHandler(handlers ...slog.Handler) *MultiHandler {
	return &MultiHandler{handlers: handlers}
}

// Enabled reports whether any of the underlying handlers process records at the given level.
func (h *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, handler := range h.handlers {
		if handler != nil && handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

// Handle dispatches the log record to all enabled underlying handlers.
// It returns the first error encountered, but handlers continue execution even if a preceding handler returns an error.
func (h *MultiHandler) Handle(ctx context.Context, record slog.Record) error {
	var firstErr error
	for _, handler := range h.handlers {
		if handler == nil {
			continue
		}
		if handler.Enabled(ctx, record.Level) {
			if err := handler.Handle(ctx, record); err != nil && firstErr == nil {
				firstErr = err
			}
		}
	}
	return firstErr
}

// WithAttrs returns a new handler containing the merged attributes applied to all configured underlying handlers.
func (h *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	handlers := make([]slog.Handler, 0, len(h.handlers))
	for _, handler := range h.handlers {
		if handler != nil {
			handlers = append(handlers, handler.WithAttrs(attrs))
		}
	}
	return NewMultiHandler(handlers...)
}

// WithGroup returns a new handler wrapping the current underlying handlers with the provided group name.
func (h *MultiHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	handlers := make([]slog.Handler, 0, len(h.handlers))
	for _, handler := range h.handlers {
		if handler != nil {
			handlers = append(handlers, handler.WithGroup(name))
		}
	}
	return NewMultiHandler(handlers...)
}
