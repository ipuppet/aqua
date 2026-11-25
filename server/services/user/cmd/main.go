package main

import (
	"aqua/pkg/infra/bootstrap"
	"aqua/services/user"
)

func main() {
	bootstrap.Run(bootstrap.Service{
		Name:    user.ServiceName,
		Factory: user.Assemble,
	})
}
