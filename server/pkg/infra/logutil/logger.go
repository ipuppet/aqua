package logutil

import (
	"errors"
	"io"
	"log/slog"
	"os"
	"sync"

	"go.opentelemetry.io/contrib/bridges/otelslog"
)

var (
	defaultLogger     *slog.Logger
	defaultLoggerOnce sync.Once
)

// Default returns a lightweight logger writing to standard output.
// It is intended for early initialization or fallback scenarios and is safe for concurrent use.
// For production use, create service-specific loggers using [NewProductionLogger].
func Default() *slog.Logger {
	// Initialize logger once using sync.Once.
	// slog.Logger is safe for concurrent use, and sync.Once guarantees
	// that defaultLogger is published safely to all goroutines.
	defaultLoggerOnce.Do(func() {
		handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelInfo,
		})
		defaultLogger = slog.New(handler)
	})
	return defaultLogger
}

// WithComponent returns a logger annotated with the specified component name.
// If logger is nil, it falls back to the package Default logger.
func WithComponent(logger *slog.Logger, component string) *slog.Logger {
	if logger == nil {
		logger = Default()
	}
	if component == "" {
		return logger
	}
	return logger.With("component", component)
}

// WithTrace returns a logger annotated with trace tracking identifiers.
// Both traceID and spanID are added if provided. Empty values are ignored.
func WithTrace(logger *slog.Logger, traceID, spanID string) *slog.Logger {
	if logger == nil {
		logger = Default()
	}

	attrs := make([]slog.Attr, 0, 2)
	if traceID != "" {
		attrs = append(attrs, slog.String("trace_id", traceID))
	}
	if spanID != "" {
		attrs = append(attrs, slog.String("span_id", spanID))
	}
	if len(attrs) == 0 {
		return logger
	}
	args := make([]any, len(attrs))
	for i, attr := range attrs {
		args[i] = attr
	}
	return logger.With(args...)
}

// NewProductionLogger returns a configured logger, a closer, and any initialization error.
// serviceName is used for file naming (e.g., user.log) and OTLP.
// The returned closer must be invoked during shutdown.
func NewProductionLogger(basePath, serviceName string) (*slog.Logger, io.Closer, error) {
	handler, closer, err := NewProductionHandler(basePath, serviceName)
	if err != nil {
		return nil, nil, err
	}
	return slog.New(handler), closer, nil
}

// NewProductionHandler returns a handler configured for production environments, writing to file and standard output.
// serviceName is used for file naming and OTLP. If OpenTelemetry is enabled, an OTLP handler is automatically included.
// The returned closer must be invoked during shutdown to ensure buffers are flushed.
func NewProductionHandler(basePath, serviceName string, extraHandlers ...slog.Handler) (slog.Handler, io.Closer, error) {
	if serviceName == "" {
		return nil, nil, errors.New("service name is required for production logger")
	}

	fileHandler, file, err := NewFileHandler(basePath, serviceName)
	if err != nil {
		return nil, nil, err
	}

	// Build handlers list: file handler + OTLP handler (if available) + extra handlers
	handlers := make([]slog.Handler, 0, 2+len(extraHandlers))
	handlers = append(handlers, fileHandler)

	if otelHandler := otelslog.NewHandler(serviceName); otelHandler != nil {
		handlers = append(handlers, otelHandler)
	}

	handlers = append(handlers, extraHandlers...)

	// Return single handler if only one, otherwise use MultiHandler
	if len(handlers) == 1 {
		return fileHandler, file, nil
	}

	return NewMultiHandler(handlers...), file, nil
}
