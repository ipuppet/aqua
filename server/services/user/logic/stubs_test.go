package logic

import (
	"context"
	"time"

	"aqua/pkg/infra/cache"
	"aqua/pkg/rbac"
	"aqua/services/user/dao"
	user "aqua/services/user/domain"
)

var _ dao.UserRepository = (*userRepoStub)(nil)
var _ dao.UserManagementRepository = (*managementRepoStub)(nil)
var _ dao.AccessControlRepository = (*managementRepoStub)(nil)

type userRepoStub struct {
	getInfoFn func(context.Context, int) (*user.UserModel, error)
	updateFn  func(context.Context, int, *user.UserUpdate) error
	cache     cache.Store
}

func (s *userRepoStub) GetUserInfo(ctx context.Context, uid int) (*user.UserModel, error) {
	if s.getInfoFn != nil {
		return s.getInfoFn(ctx, uid)
	}
	return nil, nil
}

func (s *userRepoStub) UpdateUser(ctx context.Context, uid int, updates *user.UserUpdate) error {
	if s.updateFn != nil {
		return s.updateFn(ctx, uid, updates)
	}
	return nil
}

func (s *userRepoStub) Cache() cache.Store {
	if s.cache != nil {
		return s.cache
	}
	return &mockCache{}
}

type mockCache struct{}

func (m *mockCache) Get(ctx context.Context, key string) ([]byte, error) { return nil, nil }
func (m *mockCache) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return nil
}

func (m *mockCache) Del(ctx context.Context, key string) error {
	return nil
}

func (m *mockCache) DelMatch(ctx context.Context, pattern string) error {
	return nil
}

func (m *mockCache) GetMatch(ctx context.Context, pattern string) (map[string][]byte, error) {
	return make(map[string][]byte), nil
}

func (m *mockCache) Flush(ctx context.Context) error {
	return nil
}

type managementRepoStub struct {
	userRepoStub

	searchUsersFn        func(context.Context, string, int, int) ([]user.UserModel, error)
	getUsersFn           func(context.Context, int, int) ([]user.UserModel, error)
	getUserRolesFn       func(context.Context, int) ([]rbac.Role, error)
	getRoleCodeByIDFn    func(context.Context, int) (string, error)
	getUserIDsByRoleFn   func(context.Context, int) ([]int, error)
	getUserIDsByPermFn   func(context.Context, int) ([]int, error)
	addUserRoleFn        func(context.Context, int, int) error
	deleteUserRoleFn     func(context.Context, int, int) error
	getUserPermissionsFn func(context.Context, int) ([]rbac.Permission, error)
	getRolesFn           func(context.Context) ([]rbac.Role, error)
	createRoleFn         func(context.Context, *rbac.Role) error
	updateRoleFn         func(context.Context, int, *user.RoleUpdate) error
	deleteRoleFn         func(context.Context, int) error
	getPermissionsFn     func(context.Context) ([]rbac.Permission, error)
	createPermissionFn   func(context.Context, *rbac.Permission) error
	updatePermissionFn   func(context.Context, int, *user.PermissionUpdate) error
	deletePermissionFn   func(context.Context, int) error
	getRolePermissionsFn func(context.Context, int) ([]rbac.Permission, error)
	addRolePermissionFn  func(context.Context, int, int) error
	deleteRolePermFn     func(context.Context, int, int) error
}

func (s *managementRepoStub) SearchUsers(ctx context.Context, kw string, offset, limit int) ([]user.UserModel, error) {
	if s.searchUsersFn != nil {
		return s.searchUsersFn(ctx, kw, offset, limit)
	}
	return nil, nil
}

func (s *managementRepoStub) GetUsers(ctx context.Context, offset, limit int) ([]user.UserModel, error) {
	if s.getUsersFn != nil {
		return s.getUsersFn(ctx, offset, limit)
	}
	return nil, nil
}

func (s *managementRepoStub) GetUserRoles(ctx context.Context, uid int) ([]rbac.Role, error) {
	if s.getUserRolesFn != nil {
		return s.getUserRolesFn(ctx, uid)
	}
	return nil, nil
}

func (s *managementRepoStub) GetUserIDsByRoleID(ctx context.Context, id int) ([]int, error) {
	if s.getUserIDsByRoleFn != nil {
		return s.getUserIDsByRoleFn(ctx, id)
	}
	return nil, nil
}

func (s *managementRepoStub) GetRoleCodeByID(ctx context.Context, id int) (string, error) {
	if s.getRoleCodeByIDFn != nil {
		return s.getRoleCodeByIDFn(ctx, id)
	}
	return "", nil
}

func (s *managementRepoStub) GetUserIDsByPermissionID(ctx context.Context, id int) ([]int, error) {
	if s.getUserIDsByPermFn != nil {
		return s.getUserIDsByPermFn(ctx, id)
	}
	return nil, nil
}

func (s *managementRepoStub) AddUserRole(ctx context.Context, uid int, id int) error {
	if s.addUserRoleFn != nil {
		return s.addUserRoleFn(ctx, uid, id)
	}
	return nil
}

func (s *managementRepoStub) DeleteUserRole(ctx context.Context, uid int, id int) error {
	if s.deleteUserRoleFn != nil {
		return s.deleteUserRoleFn(ctx, uid, id)
	}
	return nil
}

func (s *managementRepoStub) GetUserPermissions(ctx context.Context, uid int) ([]rbac.Permission, error) {
	if s.getUserPermissionsFn != nil {
		return s.getUserPermissionsFn(ctx, uid)
	}
	return nil, nil
}

func (s *managementRepoStub) GetRoles(ctx context.Context) ([]rbac.Role, error) {
	if s.getRolesFn != nil {
		return s.getRolesFn(ctx)
	}
	return nil, nil
}

func (s *managementRepoStub) CreateRole(ctx context.Context, role *rbac.Role) error {
	if s.createRoleFn != nil {
		return s.createRoleFn(ctx, role)
	}
	return nil
}

func (s *managementRepoStub) UpdateRole(ctx context.Context, id int, updates *user.RoleUpdate) error {
	if s.updateRoleFn != nil {
		return s.updateRoleFn(ctx, id, updates)
	}
	return nil
}

func (s *managementRepoStub) DeleteRole(ctx context.Context, id int) error {
	if s.deleteRoleFn != nil {
		return s.deleteRoleFn(ctx, id)
	}
	return nil
}

func (s *managementRepoStub) GetPermissions(ctx context.Context) ([]rbac.Permission, error) {
	if s.getPermissionsFn != nil {
		return s.getPermissionsFn(ctx)
	}
	return nil, nil
}

func (s *managementRepoStub) CreatePermission(ctx context.Context, permission *rbac.Permission) error {
	if s.createPermissionFn != nil {
		return s.createPermissionFn(ctx, permission)
	}
	return nil
}

func (s *managementRepoStub) UpdatePermission(ctx context.Context, id int, updates *user.PermissionUpdate) error {
	if s.updatePermissionFn != nil {
		return s.updatePermissionFn(ctx, id, updates)
	}
	return nil
}

func (s *managementRepoStub) DeletePermission(ctx context.Context, id int) error {
	if s.deletePermissionFn != nil {
		return s.deletePermissionFn(ctx, id)
	}
	return nil
}

func (s *managementRepoStub) GetRolePermissions(ctx context.Context, id int) ([]rbac.Permission, error) {
	if s.getRolePermissionsFn != nil {
		return s.getRolePermissionsFn(ctx, id)
	}
	return nil, nil
}

func (s *managementRepoStub) AddRolePermission(ctx context.Context, roleID int, permissionID int) error {
	if s.addRolePermissionFn != nil {
		return s.addRolePermissionFn(ctx, roleID, permissionID)
	}
	return nil
}

func (s *managementRepoStub) DeleteRolePermission(ctx context.Context, roleID int, permissionID int) error {
	if s.deleteRolePermFn != nil {
		return s.deleteRolePermFn(ctx, roleID, permissionID)
	}
	return nil
}
