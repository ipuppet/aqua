package main

import (
	"github.com/ipuppet/aqua/server/pkg/infra/bootstrap"
	"github.com/ipuppet/aqua/server/services/blog"
)

func main() {
	bootstrap.Run(bootstrap.Service{
		Name:    blog.ServiceName,
		Factory: blog.Assemble,
	})
}
