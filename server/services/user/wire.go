//go:build wireinject
// +build wireinject

package user

import (
	"github.com/google/wire"

	"github.com/ipuppet/aqua/server/pkg/infra/bootstrap"
	"github.com/ipuppet/aqua/server/pkg/infra/persist"
	"github.com/ipuppet/aqua/server/pkg/rbac"
	"github.com/ipuppet/aqua/server/pkg/transport/token"
	"github.com/ipuppet/aqua/server/services/user/dao"
	"github.com/ipuppet/aqua/server/services/user/logic"
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
