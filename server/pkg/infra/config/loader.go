package config

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Loader provides config file loading with an explicit base path (DI).
// It is responsible for path resolution; NewWatcher and callers pass relative paths.
type Loader struct {
	BasePath string
}

// NewLoader returns a Loader with the given base path.
func NewLoader(basePath string) *Loader {
	return &Loader{BasePath: basePath}
}

// resolvePath returns the absolute path for rel (relative to BasePath).
func (l *Loader) resolvePath(rel string) string {
	if filepath.IsAbs(rel) {
		return rel
	}
	return filepath.Join(l.BasePath, rel)
}

// Get loads the config file at path (relative to BasePath) into v.
func (l *Loader) Get(path string, v any) error {
	if v == nil {
		return fmt.Errorf("config destination must not be nil")
	}
	absPath := l.resolvePath(path)
	ext := strings.TrimPrefix(filepath.Ext(absPath), ".")

	cfg := viper.New()
	cfg.SetConfigFile(absPath)
	if ext != "" {
		cfg.SetConfigType(ext)
	}

	if err := cfg.ReadInConfig(); err != nil {
		return fmt.Errorf("read config %s: %w", absPath, err)
	}

	if err := cfg.Unmarshal(v); err != nil {
		return fmt.Errorf("unmarshal %s: %w", absPath, err)
	}

	return nil
}
