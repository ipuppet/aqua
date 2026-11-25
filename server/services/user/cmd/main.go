package main

import (
	"github.com/ipuppet/aqua/server/pkg/infra/bootstrap"
	"github.com/ipuppet/aqua/server/services/user"
)

func main() {
	bootstrap.Run(bootstrap.Service{
		Name:    user.ServiceName,
		Factory: user.Assemble,
	})
}
