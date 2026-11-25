package main

import (
	"github.com/ipuppet/aqua/server/pkg/infra/bootstrap"
	"github.com/ipuppet/aqua/server/services/admin"
)

func main() {
	bootstrap.Run(bootstrap.Service{
		Name:    admin.ServiceName,
		Factory: admin.Assemble,
	})
}
