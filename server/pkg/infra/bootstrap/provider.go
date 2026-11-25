package bootstrap

import (
	"github.com/google/wire"

	"aqua/pkg/infra/config"
	"aqua/pkg/infra/filestore"
	"aqua/pkg/infra/models"
)

// RuntimeBaseSet exports common runtime dependencies from Runtime.
var RuntimeBaseSet = wire.NewSet(
	wire.FieldsOf(new(*Runtime), "Env", "ConfigWatcher", "ConfigLoader", "Logger"),
	NewAppConfig,
)

// NewAppConfig returns the latest app config snapshot from the watcher.
func NewAppConfig(watcher *config.Watcher[*models.AppConfig]) *models.AppConfig {
	if watcher == nil {
		return nil
	}
	return watcher.Value()
}

// NewFileStore builds a filestore client from runtime env.
func NewFileStore(env Env) *filestore.Client {
	return filestore.NewClient(env.StorePath)
}

// FileStoreSet exports filestore path and filestore client providers together.
var FileStoreSet = wire.NewSet(
	NewFileStore,
)
