//go:build wireinject
// +build wireinject

package passport

import (
	"github.com/google/wire"

	"aqua/pkg/infra/bootstrap"
	"aqua/pkg/infra/config"
	"aqua/pkg/infra/persist"
	"aqua/pkg/rbac"
	"aqua/pkg/transport/token"
	"aqua/services/passport/dao"
	"aqua/services/passport/logic"
)

// buildModule assembles service module.
func buildModule(rt *bootstrap.Runtime) (*bootstrap.Module, error) {
	panic(wire.Build(
		bootstrap.RuntimeBaseSet,
		bootstrap.FileStoreSet,
		config.EmailConfigSet,
		config.TurnstileConfigSet,
		token.VerifierSet,
		token.PrivateKeySet,
		persist.Set,
		rbac.Set,
		dao.Set,
		logic.Set,
		provideModule,
	))
}
