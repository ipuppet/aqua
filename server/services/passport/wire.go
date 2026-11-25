//go:build wireinject
// +build wireinject

package passport

import (
	"github.com/google/wire"

	"github.com/ipuppet/aqua/server/pkg/infra/bootstrap"
	"github.com/ipuppet/aqua/server/pkg/infra/config"
	"github.com/ipuppet/aqua/server/pkg/infra/persist"
	"github.com/ipuppet/aqua/server/pkg/rbac"
	"github.com/ipuppet/aqua/server/pkg/transport/token"
	"github.com/ipuppet/aqua/server/services/passport/dao"
	"github.com/ipuppet/aqua/server/services/passport/logic"
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
