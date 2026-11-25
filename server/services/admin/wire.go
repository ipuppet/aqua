//go:build wireinject
// +build wireinject

package admin

import (
	"github.com/google/wire"

	"aqua/pkg/infra/bootstrap"
	"aqua/pkg/infra/persist"
	"aqua/pkg/rbac"
	"aqua/pkg/transport/token"
	"aqua/services/admin/dao"
	"aqua/services/admin/handler"
	"aqua/services/admin/logic"
)

// buildModule assembles service module.
func buildModule(rt *bootstrap.Runtime) (*bootstrap.Module, error) {
	panic(wire.Build(
		bootstrap.RuntimeBaseSet,
		token.VerifierSet,
		persist.Set,
		rbac.Set,
		dao.Set,
		logic.Set,
		handler.Set,
		provideModule,
	))
}
