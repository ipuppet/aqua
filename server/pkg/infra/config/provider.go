package config

import (
	"github.com/google/wire"

	"aqua/pkg/infra/models"
)

// MeiliSearchConfigProvider returns the latest MeiliSearchConfig snapshot.
type MeiliSearchConfigProvider func() models.MeiliSearchConfig

// EtcdConfigProvider returns the latest EtcdConfig snapshot.
type EtcdConfigProvider func() models.EtcdConfig

// EmailConfigProvider returns the latest EmailConfig snapshot.
type EmailConfigProvider func() models.EmailConfig

// TurnstileConfigProvider returns the latest TurnstileConfig snapshot.
type TurnstileConfigProvider func() models.TurnstileConfig

// JWTConfigProvider returns the latest JWTConfig snapshot.
type JWTConfigProvider func() models.JWTConfig

// NewMeiliSearchConfigProvider returns a runtime getter for MeiliSearchConfig.
func NewMeiliSearchConfigProvider(watcher *Watcher[*models.AppConfig]) MeiliSearchConfigProvider {
	return func() models.MeiliSearchConfig {
		if watcher == nil {
			return models.MeiliSearchConfig{}
		}
		cfg := watcher.Value()
		if cfg == nil {
			return models.MeiliSearchConfig{}
		}
		return cfg.Meili
	}
}

// MeiliSearchConfigSet exports a runtime getter for MeiliSearchConfig.
var MeiliSearchConfigSet = wire.NewSet(NewMeiliSearchConfigProvider)

// NewEtcdConfigProvider returns a runtime getter for EtcdConfig.
func NewEtcdConfigProvider(watcher *Watcher[*models.AppConfig]) EtcdConfigProvider {
	return func() models.EtcdConfig {
		if watcher == nil {
			return models.EtcdConfig{}
		}
		cfg := watcher.Value()
		if cfg == nil {
			return models.EtcdConfig{}
		}
		return cfg.Etcd
	}
}

// EtcdConfigSet exports a runtime getter for EtcdConfig.
var EtcdConfigSet = wire.NewSet(NewEtcdConfigProvider)

// NewEmailConfigProvider returns a runtime getter for EmailConfig.
func NewEmailConfigProvider(watcher *Watcher[*models.AppConfig]) EmailConfigProvider {
	return func() models.EmailConfig {
		if watcher == nil {
			return models.EmailConfig{}
		}
		cfg := watcher.Value()
		if cfg == nil {
			return models.EmailConfig{}
		}
		return cfg.Email
	}
}

// EmailConfigSet exports a runtime getter for EmailConfig.
var EmailConfigSet = wire.NewSet(NewEmailConfigProvider)

// NewTurnstileConfigProvider returns a runtime getter for TurnstileConfig.
func NewTurnstileConfigProvider(watcher *Watcher[*models.AppConfig]) TurnstileConfigProvider {
	return func() models.TurnstileConfig {
		if watcher == nil {
			return models.TurnstileConfig{}
		}
		cfg := watcher.Value()
		if cfg == nil {
			return models.TurnstileConfig{}
		}
		return cfg.Turnstile
	}
}

// TurnstileConfigSet exports a runtime getter for TurnstileConfig.
var TurnstileConfigSet = wire.NewSet(NewTurnstileConfigProvider)

// NewJWTConfigProvider returns a runtime getter for JWTConfig.
func NewJWTConfigProvider(watcher *Watcher[*models.AppConfig]) JWTConfigProvider {
	return func() models.JWTConfig {
		if watcher == nil {
			return models.JWTConfig{}
		}
		cfg := watcher.Value()
		if cfg == nil {
			return models.JWTConfig{}
		}
		return cfg.JWT
	}
}

// JWTConfigSet exports a runtime getter for JWTConfig.
var JWTConfigSet = wire.NewSet(NewJWTConfigProvider)
