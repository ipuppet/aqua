//go:build wireinject
// +build wireinject

package blog

import (
	"github.com/google/wire"

	"github.com/ipuppet/aqua/server/pkg/infra/bootstrap"
	"github.com/ipuppet/aqua/server/pkg/infra/config"
	"github.com/ipuppet/aqua/server/pkg/infra/persist"
	"github.com/ipuppet/aqua/server/pkg/rbac"
	"github.com/ipuppet/aqua/server/pkg/transport/token"
	"github.com/ipuppet/aqua/server/services/blog/dao"
	"github.com/ipuppet/aqua/server/services/blog/logic"
	"github.com/ipuppet/aqua/server/services/blog/meili"
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
