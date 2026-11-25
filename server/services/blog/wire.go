//go:build wireinject
// +build wireinject

package blog

import (
	"github.com/google/wire"

	"aqua/pkg/infra/bootstrap"
	"aqua/pkg/infra/config"
	"aqua/pkg/infra/persist"
	"aqua/pkg/rbac"
	"aqua/pkg/transport/token"
	"aqua/services/blog/dao"
	"aqua/services/blog/logic"
	"aqua/services/blog/meili"
)

// buildModule assembles service module.
func buildModule(rt *bootstrap.Runtime) (*bootstrap.Module, error) {
	panic(wire.Build(
		bootstrap.RuntimeBaseSet,
		config.MeiliSearchConfigSet,
		token.VerifierSet,
		persist.Set,
		rbac.Set,
		meili.Set,
		dao.Set,
		logic.Set,
		provideModule,
	))
}
