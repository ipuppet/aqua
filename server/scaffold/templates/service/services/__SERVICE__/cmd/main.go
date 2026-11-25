package main

import (
	"aqua/pkg/infra/bootstrap"
	"aqua/services/__SERVICE__"
)

func main() {
	bootstrap.Run(bootstrap.Service{
		Name:    __SERVICE__.ServiceName,
		Factory: __SERVICE__.Assemble,
	})
}
