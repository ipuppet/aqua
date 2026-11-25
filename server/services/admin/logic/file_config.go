package logic

import (
	"fmt"
	"log/slog"

	"github.com/ipuppet/aqua/server/pkg/infra/config"
	"github.com/ipuppet/aqua/server/services/admin/types"
)

// FileConfigProvider loads file-store settings for admin uploads.
type FileConfigProvider interface {
	FileStoreConfig() (types.FileStoreConfig, error)
}

type adminConfig struct {
	File types.FileStoreConfig `mapstructure:"file" json:"file"`
}

type fileConfigProvider struct {
	watcher *config.Watcher[adminConfig]
}

// NewFileConfigProvider constructs a provider backed by config/admin.yaml.
func NewFileConfigProvider(loader *config.Loader, logger *slog.Logger) (FileConfigProvider, error) {
	watcher, err := config.NewWatcher[adminConfig](loader, "admin.yaml", logger)
	if err != nil {
		return nil, fmt.Errorf("init admin file config provider: %w", err)
	}
	return &fileConfigProvider{watcher: watcher}, nil
}

func (p *fileConfigProvider) FileStoreConfig() (types.FileStoreConfig, error) {
	if p == nil || p.watcher == nil {
		return types.FileStoreConfig{}, fmt.Errorf("file config provider is not initialized")
	}
	return p.watcher.Value().File, nil
}
