package main

import (
	"aqua/pkg/infra/bootstrap"
	"aqua/services/admin"
)

func main() {
	bootstrap.Run(bootstrap.Service{
		Name:    admin.ServiceName,
		Factory: admin.Assemble,
	})
}
