package main

import (
	"aqua/pkg/infra/bootstrap"
	"aqua/services/blog"
)

func main() {
	bootstrap.Run(bootstrap.Service{
		Name:    blog.ServiceName,
		Factory: blog.Assemble,
	})
}
