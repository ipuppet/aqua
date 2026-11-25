package main

import (
	"github.com/ipuppet/aqua/server/pkg/infra/bootstrap"
	"github.com/ipuppet/aqua/server/services/passport"
)

func main() {
	bootstrap.Run(bootstrap.Service{
		Name:    passport.ServiceName,
		Factory: passport.Assemble,
	})
}
