package dao

import (
	"context"

	"github.com/ipuppet/aqua/server/pkg/infra/cache"
	"github.com/ipuppet/aqua/server/pkg/infra/database"
	"github.com/ipuppet/aqua/server/pkg/rbac"
	user "github.com/ipuppet/aqua/server/services/user/domain"
)

// UserRepository exposes read/write operations required by public and private handlers.
type UserRepository interface {
	GetUserInfo(ctx context.Context, uid int) (*user.UserModel, error)
	UpdateUser(ctx context.Context, uid int, updates *user.UserUpdate) error
	Cache() cache.Store
}

// UserManagementRepository provides admin-side user list/search/update persistence operations.
type UserManagementRepository interface {
	Cache() cache.Store
	SearchUsers(ctx context.Context, kw string, offset, limit int) ([]user.UserModel, error)
	GetUsers(ctx context.Context, offset, limit int) ([]user.UserModel, error)
	UpdateUser(ctx context.Context, uid int, updates *user.UserUpdate) error
}

// AccessControlRepository provides RBAC and user-role relationship persistence operations.
type AccessControlRepository interface {
	Cache() cache.Store
	GetRoleCodeByID(ctx context.Context, id int) (string, error)
	GetUserIDsByRoleID(ctx context.Context, id int) ([]int, error)
	GetUserIDsByPermissionID(ctx context.Context, id int) ([]int, error)
	GetUserRoles(ctx context.Context, uid int) ([]rbac.Role, error)
	AddUserRole(ctx context.Context, uid int, id int) error
	DeleteUserRole(ctx context.Context, uid int, id int) error
	GetUserPermissions(ctx context.Context, uid int) ([]rbac.Permission, error)
	GetRoles(ctx context.Context) ([]rbac.Role, error)
	CreateRole(ctx context.Context, role *rbac.Role) error
	UpdateRole(ctx context.Context, id int, updates *user.RoleUpdate) error
	DeleteRole(ctx context.Context, id int) error
	GetPermissions(ctx context.Context) ([]rbac.Permission, error)
	CreatePermission(ctx context.Context, permission *rbac.Permission) error
	UpdatePermission(ctx context.Context, id int, updates *user.PermissionUpdate) error
	DeletePermission(ctx context.Context, id int) error
	GetRolePermissions(ctx context.Context, id int) ([]rbac.Permission, error)
	AddRolePermission(ctx context.Context, roleID int, permissionID int) error
	DeleteRolePermission(ctx context.Context, roleID int, permissionID int) error
}

// repository implements all user service persistence ports.
type repository struct {
	db         *database.Client
	cache      cache.Store
	rbacReader rbac.Reader
}

// NewRepository builds a repository backed by the provided database connection.
func NewRepository(db *database.Client, c cache.Store, rbacReader rbac.Reader) *repository {
	return &repository{db: db, cache: c, rbacReader: rbacReader}
}

func (r *repository) Cache() cache.Store {
	return r.cache
}

var _ UserRepository = (*repository)(nil)
var _ UserManagementRepository = (*repository)(nil)
var _ AccessControlRepository = (*repository)(nil)
