package main

import (
	"aqua/pkg/infra/bootstrap"
	"aqua/services/admin"
	"aqua/services/blog"
	"aqua/services/passport"
	"aqua/services/user"
)

func main() {
	bootstrap.RunAll([]bootstrap.Service{
		{Name: admin.ServiceName, Factory: admin.Assemble},
		{Name: passport.ServiceName, Factory: passport.Assemble},
		{Name: user.ServiceName, Factory: user.Assemble},
		{Name: blog.ServiceName, Factory: blog.Assemble},
	})
}
