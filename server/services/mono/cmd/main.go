package main

import (
	"github.com/ipuppet/aqua/server/pkg/infra/bootstrap"
	"github.com/ipuppet/aqua/server/services/admin"
	"github.com/ipuppet/aqua/server/services/blog"
	"github.com/ipuppet/aqua/server/services/passport"
	"github.com/ipuppet/aqua/server/services/user"
)

func main() {
	bootstrap.RunAll([]bootstrap.Service{
		{Name: admin.ServiceName, Factory: admin.Assemble},
		{Name: passport.ServiceName, Factory: passport.Assemble},
		{Name: user.ServiceName, Factory: user.Assemble},
		{Name: blog.ServiceName, Factory: blog.Assemble},
	})
}
