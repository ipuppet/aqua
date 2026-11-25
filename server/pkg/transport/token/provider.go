package token

import (
	"crypto/rsa"
	"fmt"

	"github.com/google/wire"

	"aqua/pkg/infra/config"
)

// PublicKeySet provides JWT public key for token verification.
var PublicKeySet = wire.NewSet(NewPublicKey)

// PrivateKeySet provides JWT private key for token signing.
var PrivateKeySet = wire.NewSet(NewPrivateKey)

// VerifierSet provides JWT config and public key for token verification.
var VerifierSet = wire.NewSet(
	config.JWTConfigSet,
	PublicKeySet,
)

// NewPublicKey loads RSA public key from JWT config path.
func NewPublicKey(configProvider config.JWTConfigProvider) (*rsa.PublicKey, error) {
	if configProvider == nil {
		return nil, fmt.Errorf("jwt config provider is not configured")
	}
	jwtCfg := configProvider()
	return LoadPublicKey(jwtCfg.PublicKeyPath)
}

// NewPrivateKey loads RSA private key from JWT config path.
func NewPrivateKey(configProvider config.JWTConfigProvider) (*rsa.PrivateKey, error) {
	if configProvider == nil {
		return nil, fmt.Errorf("jwt config provider is not configured")
	}
	jwtCfg := configProvider()
	return LoadPrivateKey(jwtCfg.PrivateKeyPath)
}
