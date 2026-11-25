package logutil

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

const (
	// logDirPerm defines the standard permissions for log directories.
	logDirPerm = 0o755
	// logFilePerm defines the standard permissions for individual log files.
	logFilePerm = 0o644
)

// handlerConfig represents the configuration for constructing a slog handler.
type handlerConfig struct {
	Level     slog.Leveler
	AddSource bool
	Writers   []io.Writer
}

// NewFileHandler creates a handler that writes logs to a local file and standard output.
// serviceName is used for the log file name (e.g. user.log). It creates the log directory if needed
// and opens the log file in append mode. The returned closer must be invoked during shutdown.
func NewFileHandler(basePath, serviceName string) (slog.Handler, io.Closer, error) {
	if basePath == "" {
		handler := newHandler(handlerConfig{
			Level:     slog.LevelInfo,
			AddSource: true,
			Writers:   []io.Writer{os.Stdout},
		})
		return handler, nil, nil
	}

	logFile := filepath.Join(basePath, serviceName+".log")
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, logFilePerm)
	if err != nil {
		return nil, nil, err
	}

	handler := newHandler(handlerConfig{
		Level:     slog.LevelInfo,
		AddSource: true,
		Writers:   []io.Writer{file, os.Stdout},
	})
	return handler, file, nil
}

// newHandler returns a configured JSON handler based on the provided configuration.
func newHandler(cfg handlerConfig) slog.Handler {
	opts := &slog.HandlerOptions{
		AddSource: cfg.AddSource,
	}
	if cfg.Level != nil {
		opts.Level = cfg.Level
	}

	writer := buildWriter(cfg.Writers)
	return slog.NewJSONHandler(writer, opts)
}

// buildWriter returns a customized io.Writer from multiple writers.
// It filters out nil writers and returns standard output if no valid writers are present.
func buildWriter(writers []io.Writer) io.Writer {
	if len(writers) == 0 {
		return os.Stdout
	}

	valid := make([]io.Writer, 0, len(writers))
	for _, w := range writers {
		if w != nil {
			valid = append(valid, w)
		}
	}

	if len(valid) == 0 {
		return os.Stdout
	}
	if len(valid) == 1 {
		return valid[0]
	}
	return io.MultiWriter(valid...)
}
