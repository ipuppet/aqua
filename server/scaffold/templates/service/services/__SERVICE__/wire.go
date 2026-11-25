//go:build wireinject
// +build wireinject

package __SERVICE__

import (
	"github.com/google/wire"

	"aqua/pkg/infra/bootstrap"
	"aqua/services/__SERVICE__/logic"
)

// buildModule assembles service module.
func buildModule(rt *bootstrap.Runtime) (*bootstrap.Module, error) {
	panic(wire.Build(
		bootstrap.RuntimeBaseSet,
		logic.Set,
		provideModule,
	))
}
