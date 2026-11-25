//go:build wireinject
// +build wireinject

package user

import (
	"github.com/google/wire"

	"aqua/pkg/infra/bootstrap"
	"aqua/pkg/infra/persist"
	"aqua/pkg/rbac"
	"aqua/pkg/transport/token"
	"aqua/services/user/dao"
	"aqua/services/user/logic"
)

// buildModule assembles service module.
func buildModule(rt *bootstrap.Runtime) (*bootstrap.Module, error) {
	panic(wire.Build(
		bootstrap.RuntimeBaseSet,
		bootstrap.FileStoreSet,
		token.VerifierSet,
		persist.Set,
		rbac.Set,
		rbac.EventBusSet,
		dao.Set,
		logic.Set,
		provideModule,
	))
}
