package main

import (
	"aqua/pkg/infra/bootstrap"
	"aqua/services/passport"
)

func main() {
	bootstrap.Run(bootstrap.Service{
		Name:    passport.ServiceName,
		Factory: passport.Assemble,
	})
}
