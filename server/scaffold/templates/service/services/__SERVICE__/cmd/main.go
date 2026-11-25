package main

import (
	"github.com/ipuppet/aqua/server/pkg/infra/bootstrap"
	"github.com/ipuppet/aqua/server/services/__SERVICE__"
)

func main() {
	bootstrap.Run(bootstrap.Service{
		Name:    __SERVICE__.ServiceName,
		Factory: __SERVICE__.Assemble,
	})
}
