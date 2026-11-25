package grpcserver

import (
	"context"
	"database/sql"

	userv1 "aqua/gen/go/aqua/user/v1"
	"aqua/pkg/rbac"
	"aqua/pkg/types"
	user "aqua/services/user/domain"
	"aqua/services/user/logic"
)

// AdminServer provides gRPC APIs for admin-side user operations.
type AdminServer struct {
	userManagementService *logic.UserManagementService
	accessControlService  *logic.AccessControlService
	userv1.UnimplementedAdminUserServiceServer
}

// NewAdminServer constructs AdminServer.
func NewAdminServer(userManagementService *logic.UserManagementService, accessControlService *logic.AccessControlService) *AdminServer {
	return &AdminServer{userManagementService: userManagementService, accessControlService: accessControlService}
}

// --- User Management ---

func (s *AdminServer) ListUsers(ctx context.Context, req *userv1.ListUsersRequest) (*userv1.ListUsersResponse, error) {
	items, err := s.userManagementService.GetUsers(ctx, int(req.GetPage()), int(req.GetPageSize()))
	if err != nil {
		return nil, err
	}
	return &userv1.ListUsersResponse{Users: ToPBUsers(items)}, nil
}

func (s *AdminServer) SearchUsers(ctx context.Context, req *userv1.SearchUsersRequest) (*userv1.SearchUsersResponse, error) {
	items, err := s.userManagementService.SearchUsers(ctx, req.GetKeyword(), int(req.GetPage()), int(req.GetPageSize()))
	if err != nil {
		return nil, err
	}
	return &userv1.SearchUsersResponse{Users: ToPBUsers(items)}, nil
}

func (s *AdminServer) UpdateUser(ctx context.Context, req *userv1.UpdateUserRequest) (*userv1.UpdateUserResponse, error) {
	updates := &user.UserUpdate{}
	if req.Nickname != nil {
		v := req.Nickname.GetValue()
		updates.Nickname = &v
	}
	if req.Email != nil {
		v := req.Email.GetValue()
		updates.Email = &v
	}
	if req.Phone != nil {
		v := req.Phone.GetValue()
		updates.Phone = &v
	}
	if req.Avatar != nil {
		v := req.Avatar.GetValue()
		updates.Avatar = &v
	}
	if req.Gender != nil {
		v := int(req.Gender.GetValue())
		updates.Gender = &v
	}
	if req.Lock != nil {
		v := int(req.Lock.GetValue())
		updates.Lock = &v
	}
	if err := s.userManagementService.UpdateUser(ctx, int(req.GetUid()), updates); err != nil {
		return nil, err
	}
	return &userv1.UpdateUserResponse{}, nil
}

// --- RBAC ---

func (s *AdminServer) ListUserRoles(ctx context.Context, req *userv1.ListUserRolesRequest) (*userv1.ListUserRolesResponse, error) {
	items, err := s.accessControlService.GetUserRoles(ctx, int(req.GetUid()))
	if err != nil {
		return nil, err
	}
	return &userv1.ListUserRolesResponse{Roles: ToPBRoles(items)}, nil
}

func (s *AdminServer) AddUserRole(ctx context.Context, req *userv1.AddUserRoleRequest) (*userv1.AddUserRoleResponse, error) {
	if err := s.accessControlService.AddUserRole(ctx, int(req.GetUid()), int(req.GetId())); err != nil {
		return nil, err
	}
	return &userv1.AddUserRoleResponse{}, nil
}

func (s *AdminServer) DeleteUserRole(ctx context.Context, req *userv1.DeleteUserRoleRequest) (*userv1.DeleteUserRoleResponse, error) {
	if err := s.accessControlService.DeleteUserRole(ctx, int(req.GetUid()), int(req.GetId())); err != nil {
		return nil, err
	}
	return &userv1.DeleteUserRoleResponse{}, nil
}

func (s *AdminServer) ListUserPermissions(ctx context.Context, req *userv1.ListUserPermissionsRequest) (*userv1.ListUserPermissionsResponse, error) {
	items, err := s.accessControlService.GetUserPermissions(ctx, int(req.GetUid()))
	if err != nil {
		return nil, err
	}
	return &userv1.ListUserPermissionsResponse{Permissions: ToPBPermissions(items)}, nil
}

func (s *AdminServer) ListRoles(ctx context.Context, _ *userv1.ListRolesRequest) (*userv1.ListRolesResponse, error) {
	items, err := s.accessControlService.GetRoles(ctx)
	if err != nil {
		return nil, err
	}
	return &userv1.ListRolesResponse{Roles: ToPBRoles(items)}, nil
}

func (s *AdminServer) CreateRole(ctx context.Context, req *userv1.CreateRoleRequest) (*userv1.CreateRoleResponse, error) {
	role := &rbac.Role{
		Code: req.GetCode(),
		Description: types.NullString{
			NullString: sql.NullString{String: req.GetDescription(), Valid: req.GetDescription() != ""},
		},
	}
	if err := s.accessControlService.CreateRole(ctx, role); err != nil {
		return nil, err
	}
	return &userv1.CreateRoleResponse{}, nil
}

func (s *AdminServer) UpdateRole(ctx context.Context, req *userv1.UpdateRoleRequest) (*userv1.UpdateRoleResponse, error) {
	code := req.GetCode()
	desc := req.GetDescription()
	updates := &user.RoleUpdate{
		Code:        &code,
		Description: &desc,
	}
	if err := s.accessControlService.UpdateRole(ctx, int(req.GetId()), updates); err != nil {
		return nil, err
	}
	return &userv1.UpdateRoleResponse{}, nil
}

func (s *AdminServer) DeleteRole(ctx context.Context, req *userv1.DeleteRoleRequest) (*userv1.DeleteRoleResponse, error) {
	if err := s.accessControlService.DeleteRole(ctx, int(req.GetId())); err != nil {
		return nil, err
	}
	return &userv1.DeleteRoleResponse{}, nil
}

func (s *AdminServer) ListPermissions(ctx context.Context, _ *userv1.ListPermissionsRequest) (*userv1.ListPermissionsResponse, error) {
	items, err := s.accessControlService.GetPermissions(ctx)
	if err != nil {
		return nil, err
	}
	return &userv1.ListPermissionsResponse{Permissions: ToPBPermissions(items)}, nil
}

func (s *AdminServer) CreatePermission(ctx context.Context, req *userv1.CreatePermissionRequest) (*userv1.CreatePermissionResponse, error) {
	permission := &rbac.Permission{
		Code: req.GetCode(),
		Description: types.NullString{
			NullString: sql.NullString{String: req.GetDescription(), Valid: req.GetDescription() != ""},
		},
	}
	if err := s.accessControlService.CreatePermission(ctx, permission); err != nil {
		return nil, err
	}
	return &userv1.CreatePermissionResponse{}, nil
}

func (s *AdminServer) UpdatePermission(ctx context.Context, req *userv1.UpdatePermissionRequest) (*userv1.UpdatePermissionResponse, error) {
	code := req.GetCode()
	desc := req.GetDescription()
	updates := &user.PermissionUpdate{
		Code:        &code,
		Description: &desc,
	}
	if err := s.accessControlService.UpdatePermission(ctx, int(req.GetId()), updates); err != nil {
		return nil, err
	}
	return &userv1.UpdatePermissionResponse{}, nil
}

func (s *AdminServer) DeletePermission(ctx context.Context, req *userv1.DeletePermissionRequest) (*userv1.DeletePermissionResponse, error) {
	if err := s.accessControlService.DeletePermission(ctx, int(req.GetId())); err != nil {
		return nil, err
	}
	return &userv1.DeletePermissionResponse{}, nil
}

func (s *AdminServer) ListRolePermissions(ctx context.Context, req *userv1.ListRolePermissionsRequest) (*userv1.ListRolePermissionsResponse, error) {
	items, err := s.accessControlService.GetRolePermissions(ctx, int(req.GetId()))
	if err != nil {
		return nil, err
	}
	return &userv1.ListRolePermissionsResponse{Permissions: ToPBPermissions(items)}, nil
}

func (s *AdminServer) AddRolePermission(ctx context.Context, req *userv1.AddRolePermissionRequest) (*userv1.AddRolePermissionResponse, error) {
	if err := s.accessControlService.AddRolePermission(ctx, int(req.GetRoleId()), int(req.GetPermissionId())); err != nil {
		return nil, err
	}
	return &userv1.AddRolePermissionResponse{}, nil
}

func (s *AdminServer) DeleteRolePermission(ctx context.Context, req *userv1.DeleteRolePermissionRequest) (*userv1.DeleteRolePermissionResponse, error) {
	if err := s.accessControlService.DeleteRolePermission(ctx, int(req.GetRoleId()), int(req.GetPermissionId())); err != nil {
		return nil, err
	}
	return &userv1.DeleteRolePermissionResponse{}, nil
}
