package logic

import (
	"context"
	"log/slog"

	"github.com/ipuppet/aqua/server/pkg/infra/logutil"
	"github.com/ipuppet/aqua/server/pkg/rbac"
	"github.com/ipuppet/aqua/server/pkg/transport/token"
	"github.com/ipuppet/aqua/server/services/user/dao"
	user "github.com/ipuppet/aqua/server/services/user/domain"
)

type AccessControlService struct {
	repo dao.AccessControlRepository
	bus  rbac.EventBus
}

// NewAccessControlService constructs AccessControlService.
func NewAccessControlService(repo dao.AccessControlRepository, bus rbac.EventBus) *AccessControlService {
	if bus == nil {
		handler := rbac.NewCacheInvalidationHandler(repo.Cache())
		bus = rbac.NewEventBus(handler)
	}
	return &AccessControlService{repo: repo, bus: bus}
}

func (s *AccessControlService) log(ctx context.Context) *slog.Logger {
	return logutil.FromContext(ctx)
}

func (s *AccessControlService) AddUserRole(ctx context.Context, uid int, id int) error {
	if err := s.repo.AddUserRole(ctx, uid, id); err != nil {
		s.log(ctx).Error("failed to add role to user",
			slog.Int("user.uid", uid),
			slog.Int("role.id", id),
			slog.Any("error", err))
		return err
	}
	s.bus.Publish(ctx, rbac.UserRoleAssigned{UID: uid})
	if err := token.RevokeUserRefreshTokens(ctx, s.repo.Cache(), uid); err != nil {
		s.log(ctx).Warn("failed to revoke user refresh tokens after role assignment",
			slog.Int("user.uid", uid),
			slog.Int("role.id", id),
			slog.Any("error", err))
	}
	return nil
}

func (s *AccessControlService) DeleteUserRole(ctx context.Context, uid int, id int) error {
	if err := s.repo.DeleteUserRole(ctx, uid, id); err != nil {
		s.log(ctx).Error("failed to remove role from user",
			slog.Int("user.uid", uid),
			slog.Int("role.id", id),
			slog.Any("error", err))
		return err
	}
	s.bus.Publish(ctx, rbac.UserRoleRevoked{UID: uid})
	if err := token.RevokeUserRefreshTokens(ctx, s.repo.Cache(), uid); err != nil {
		s.log(ctx).Warn("failed to revoke user refresh tokens after role revocation",
			slog.Int("user.uid", uid),
			slog.Int("role.id", id),
			slog.Any("error", err))
	}
	return nil
}

func (s *AccessControlService) GetUserRoles(ctx context.Context, uid int) ([]user.RoleResponse, error) {
	roles, err := s.repo.GetUserRoles(ctx, uid)
	if err != nil {
		s.log(ctx).Error("failed to get user roles",
			slog.Int("user.uid", uid),
			slog.Any("error", err))
		return nil, err
	}
	responses := make([]user.RoleResponse, 0, len(roles))
	for _, role := range roles {
		responses = append(responses, *ToRoleResponse(&role))
	}
	return responses, nil
}

func (s *AccessControlService) GetUserPermissions(ctx context.Context, uid int) ([]user.PermissionResponse, error) {
	permissions, err := s.repo.GetUserPermissions(ctx, uid)
	if err != nil {
		s.log(ctx).Error("failed to get user permissions",
			slog.Int("user.uid", uid),
			slog.Any("error", err))
		return nil, err
	}
	responses := make([]user.PermissionResponse, 0, len(permissions))
	for _, permission := range permissions {
		responses = append(responses, *ToPermissionResponse(&permission))
	}
	return responses, nil
}

func (s *AccessControlService) GetRoles(ctx context.Context) ([]user.RoleResponse, error) {
	roles, err := s.repo.GetRoles(ctx)
	if err != nil {
		s.log(ctx).Error("failed to get all roles",
			slog.Any("error", err))
		return nil, err
	}
	responses := make([]user.RoleResponse, 0, len(roles))
	for _, role := range roles {
		responses = append(responses, *ToRoleResponse(&role))
	}
	return responses, nil
}

func (s *AccessControlService) CreateRole(ctx context.Context, role *rbac.Role) error {
	if err := s.repo.CreateRole(ctx, role); err != nil {
		s.log(ctx).Error("failed to create role",
			slog.String("role.code", role.Code),
			slog.Any("error", err))
		return err
	}
	return nil
}

func (s *AccessControlService) UpdateRole(ctx context.Context, id int, updates *user.RoleUpdate) error {
	roleCode, err := s.repo.GetRoleCodeByID(ctx, id)
	if err != nil {
		s.log(ctx).Error("failed to resolve role code by role id",
			slog.Int("role.id", id),
			slog.Any("error", err))
		return err
	}

	affectedUids, err := s.repo.GetUserIDsByRoleID(ctx, id)
	if err != nil {
		s.log(ctx).Error("failed to resolve users affected by role update",
			slog.Int("role.id", id),
			slog.Any("error", err))
		return err
	}
	if err := s.repo.UpdateRole(ctx, id, updates); err != nil {
		s.log(ctx).Error("failed to update role",
			slog.Int("role.id", id),
			slog.Any("error", err))
		return err
	}
	s.bus.Publish(ctx, rbac.RoleUpdated{Role: roleCode})
	s.revokeUsersRefreshTokens(ctx, affectedUids, "role updated",
		slog.Int("role.id", id))
	return nil
}

func (s *AccessControlService) DeleteRole(ctx context.Context, id int) error {
	roleCode, err := s.repo.GetRoleCodeByID(ctx, id)
	if err != nil {
		s.log(ctx).Error("failed to resolve role code by role id",
			slog.Int("role.id", id),
			slog.Any("error", err))
		return err
	}

	affectedUids, err := s.repo.GetUserIDsByRoleID(ctx, id)
	if err != nil {
		s.log(ctx).Error("failed to resolve users affected by role deletion",
			slog.Int("role.id", id),
			slog.Any("error", err))
		return err
	}
	if err := s.repo.DeleteRole(ctx, id); err != nil {
		s.log(ctx).Error("failed to delete role",
			slog.Int("role.id", id),
			slog.Any("error", err))
		return err
	}
	s.bus.Publish(ctx, rbac.RoleDeleted{Role: roleCode})
	s.revokeUsersRefreshTokens(ctx, affectedUids, "role deleted",
		slog.Int("role.id", id))
	return nil
}

func (s *AccessControlService) GetPermissions(ctx context.Context) ([]user.PermissionResponse, error) {
	permissions, err := s.repo.GetPermissions(ctx)
	if err != nil {
		s.log(ctx).Error("failed to get all permissions",
			slog.Any("error", err))
		return nil, err
	}
	responses := make([]user.PermissionResponse, 0, len(permissions))
	for _, permission := range permissions {
		responses = append(responses, *ToPermissionResponse(&permission))
	}
	return responses, nil
}

func (s *AccessControlService) CreatePermission(ctx context.Context, permission *rbac.Permission) error {
	if err := s.repo.CreatePermission(ctx, permission); err != nil {
		s.log(ctx).Error("failed to create permission",
			slog.String("permission.code", permission.Code),
			slog.Any("error", err))
		return err
	}
	return nil
}

func (s *AccessControlService) UpdatePermission(ctx context.Context, id int, updates *user.PermissionUpdate) error {
	affectedUids, err := s.repo.GetUserIDsByPermissionID(ctx, id)
	if err != nil {
		s.log(ctx).Error("failed to resolve users affected by permission update",
			slog.Int("permission.id", id),
			slog.Any("error", err))
		return err
	}
	if err := s.repo.UpdatePermission(ctx, id, updates); err != nil {
		s.log(ctx).Error("failed to update permission",
			slog.Int("permission.id", id),
			slog.Any("error", err))
		return err
	}
	s.bus.Publish(ctx, rbac.PermissionUpdated{})
	s.revokeUsersRefreshTokens(ctx, affectedUids, "permission updated",
		slog.Int("permission.id", id))
	return nil
}

func (s *AccessControlService) DeletePermission(ctx context.Context, id int) error {
	affectedUids, err := s.repo.GetUserIDsByPermissionID(ctx, id)
	if err != nil {
		s.log(ctx).Error("failed to resolve users affected by permission deletion",
			slog.Int("permission.id", id),
			slog.Any("error", err))
		return err
	}
	if err := s.repo.DeletePermission(ctx, id); err != nil {
		s.log(ctx).Error("failed to delete permission",
			slog.Int("permission.id", id),
			slog.Any("error", err))
		return err
	}
	s.bus.Publish(ctx, rbac.PermissionDeleted{})
	s.revokeUsersRefreshTokens(ctx, affectedUids, "permission deleted",
		slog.Int("permission.id", id))
	return nil
}

func (s *AccessControlService) GetRolePermissions(ctx context.Context, id int) ([]user.PermissionResponse, error) {
	permissions, err := s.repo.GetRolePermissions(ctx, id)
	if err != nil {
		s.log(ctx).Error("failed to get role permissions",
			slog.Int("role.id", id),
			slog.Any("error", err))
		return nil, err
	}
	responses := make([]user.PermissionResponse, 0, len(permissions))
	for _, permission := range permissions {
		responses = append(responses, *ToPermissionResponse(&permission))
	}
	return responses, nil
}

func (s *AccessControlService) AddRolePermission(ctx context.Context, roleID int, permissionID int) error {
	roleCode, err := s.repo.GetRoleCodeByID(ctx, roleID)
	if err != nil {
		s.log(ctx).Error("failed to resolve role code by role id",
			slog.Int("role.id", roleID),
			slog.Int("permission.id", permissionID),
			slog.Any("error", err))
		return err
	}

	if err := s.repo.AddRolePermission(ctx, roleID, permissionID); err != nil {
		s.log(ctx).Error("failed to add permission to role",
			slog.Int("role.id", roleID),
			slog.Int("permission.id", permissionID),
			slog.Any("error", err))
		return err
	}
	affectedUids, err := s.repo.GetUserIDsByRoleID(ctx, roleID)
	if err != nil {
		s.log(ctx).Error("failed to resolve users affected by role permission assignment",
			slog.Int("role.id", roleID),
			slog.Int("permission.id", permissionID),
			slog.Any("error", err))
		return err
	}
	s.bus.Publish(ctx, rbac.RolePermissionAssigned{Role: roleCode})
	s.revokeUsersRefreshTokens(ctx, affectedUids, "role permission assigned",
		slog.Int("role.id", roleID),
		slog.Int("permission.id", permissionID))
	return nil
}

func (s *AccessControlService) DeleteRolePermission(ctx context.Context, roleID int, permissionID int) error {
	roleCode, err := s.repo.GetRoleCodeByID(ctx, roleID)
	if err != nil {
		s.log(ctx).Error("failed to resolve role code by role id",
			slog.Int("role.id", roleID),
			slog.Int("permission.id", permissionID),
			slog.Any("error", err))
		return err
	}

	if err := s.repo.DeleteRolePermission(ctx, roleID, permissionID); err != nil {
		s.log(ctx).Error("failed to remove permission from role",
			slog.Int("role.id", roleID),
			slog.Int("permission.id", permissionID),
			slog.Any("error", err))
		return err
	}
	affectedUids, err := s.repo.GetUserIDsByRoleID(ctx, roleID)
	if err != nil {
		s.log(ctx).Error("failed to resolve users affected by role permission revocation",
			slog.Int("role.id", roleID),
			slog.Int("permission.id", permissionID),
			slog.Any("error", err))
		return err
	}
	s.bus.Publish(ctx, rbac.RolePermissionRevoked{Role: roleCode})
	s.revokeUsersRefreshTokens(ctx, affectedUids, "role permission revoked",
		slog.Int("role.id", roleID),
		slog.Int("permission.id", permissionID))
	return nil
}

func (s *AccessControlService) revokeUsersRefreshTokens(ctx context.Context, uids []int, reason string, attrs ...slog.Attr) {
	for _, uid := range uids {
		if err := token.RevokeUserRefreshTokens(ctx, s.repo.Cache(), uid); err != nil {
			logger := s.log(ctx).With(
				slog.String("reason", reason),
				slog.Int("user.uid", uid),
			)
			for _, attr := range attrs {
				logger = logger.With(attr)
			}
			logger.Warn("failed to revoke user refresh tokens", slog.Any("error", err))
		}
	}
}
