//go:build wireinject
// +build wireinject

package admin

import (
	"github.com/google/wire"

	"github.com/ipuppet/aqua/server/pkg/infra/bootstrap"
	"github.com/ipuppet/aqua/server/pkg/infra/persist"
	"github.com/ipuppet/aqua/server/pkg/rbac"
	"github.com/ipuppet/aqua/server/pkg/transport/token"
	"github.com/ipuppet/aqua/server/services/admin/dao"
	"github.com/ipuppet/aqua/server/services/admin/handler"
	"github.com/ipuppet/aqua/server/services/admin/logic"
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
