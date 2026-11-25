package dao

import (
	"context"
	"fmt"

	"github.com/ipuppet/aqua/server/pkg/infra/database"
	"github.com/ipuppet/aqua/server/pkg/rbac"
	user "github.com/ipuppet/aqua/server/services/user/domain"
)

func (r *repository) GetUserIDsByRoleID(ctx context.Context, id int) ([]int, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	uids := []int{}
	if err := database.Retrieve(ctx, r.db, &uids,
		`SELECT DISTINCT uid FROM ums.rbac_user_role WHERE role_id = $1 ORDER BY uid`, id); err != nil {
		return nil, fmt.Errorf("retrieving user ids by role id %d: %w", id, err)
	}
	return uids, nil
}

func (r *repository) GetRoleCodeByID(ctx context.Context, id int) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}

	var roleCode string
	if err := r.db.QueryRowContext(ctx, `SELECT code FROM ums.rbac_role WHERE id = $1`, id).Scan(&roleCode); err != nil {
		return "", fmt.Errorf("retrieving role code by role id %d: %w", id, err)
	}
	return roleCode, nil
}

// GetUserIDsByPermissionID returns the list of user IDs that have the specified permission through their assigned roles.
// This is used for cache invalidation when permission-role relationships are updated.
func (r *repository) GetUserIDsByPermissionID(ctx context.Context, id int) ([]int, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	uids := []int{}
	if err := database.Retrieve(ctx, r.db, &uids,
		`SELECT DISTINCT ur.uid
		FROM ums.rbac_user_role ur
		JOIN ums.rbac_role_permission rp ON ur.role_id = rp.role_id
		WHERE rp.permission_id = $1
		ORDER BY ur.uid`, id); err != nil {
		return nil, fmt.Errorf("retrieving user ids by permission id %d: %w", id, err)
	}
	return uids, nil
}

func (r *repository) AddUserRole(ctx context.Context, uid int, id int) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if _, err := r.db.ExecContext(ctx, "insert into ums.rbac_user_role (uid, role_id) values ($1, $2)", uid, id); err != nil {
		return fmt.Errorf("adding id %d to user_id %d: %w", id, uid, err)
	}
	return nil
}

func (r *repository) DeleteUserRole(ctx context.Context, uid int, id int) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if _, err := r.db.ExecContext(ctx, "delete from ums.rbac_user_role where uid=$1 and role_id=$2", uid, id); err != nil {
		return fmt.Errorf("deleting id %d from user_id %d: %w", id, uid, err)
	}
	return nil
}

func (r *repository) GetUserRoles(ctx context.Context, uid int) ([]rbac.Role, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	roles, err := r.rbacReader.LoadUserRoles(ctx, uid)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (r *repository) GetUserPermissions(ctx context.Context, uid int) ([]rbac.Permission, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	permissions, err := r.rbacReader.LoadUserPermissions(ctx, uid)
	if err != nil {
		return nil, err
	}

	return permissions, nil
}

func (r *repository) GetRoles(ctx context.Context) ([]rbac.Role, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var roles []rbac.Role
	if err := database.Retrieve(ctx, r.db, &roles, "select id, code, description from ums.rbac_role order by id"); err != nil {
		return nil, fmt.Errorf("retrieving roles: %w", err)
	}

	return roles, nil
}

func (r *repository) CreateRole(ctx context.Context, role *rbac.Role) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if _, err := database.MustExec(r.db.ExecContext(ctx, `insert into ums.rbac_role
		(code, description)
		values
		($1, $2)`,
		role.Code, role.Description)); err != nil {
		return fmt.Errorf("adding role code %s: %w", role.Code, err)
	}
	return nil
}

func (r *repository) UpdateRole(ctx context.Context, id int, updates *user.RoleUpdate) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if err := database.ExecuteUpdate(r.db, "ums.rbac_role", updates, "id=$1", id); err != nil {
		return fmt.Errorf("updating role id %d: %w", id, err)
	}
	return nil
}

func (r *repository) DeleteRole(ctx context.Context, id int) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if _, err := database.MustExec(r.db.ExecContext(ctx, "delete from ums.rbac_role where id=$1", id)); err != nil {
		return fmt.Errorf("deleting role id %d: %w", id, err)
	}
	return nil
}

func (r *repository) GetPermissions(ctx context.Context) ([]rbac.Permission, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var permissions []rbac.Permission
	if err := database.Retrieve(ctx, r.db, &permissions, "select id, code, description from ums.rbac_permission order by id"); err != nil {
		return nil, fmt.Errorf("retrieving permissions: %w", err)
	}

	return permissions, nil
}

func (r *repository) CreatePermission(ctx context.Context, permission *rbac.Permission) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if _, err := database.MustExec(r.db.ExecContext(ctx, `insert into ums.rbac_permission
		(code, description)
		values
		($1, $2)`,
		permission.Code, permission.Description)); err != nil {
		return fmt.Errorf("adding permission code %s: %w", permission.Code, err)
	}
	return nil
}

func (r *repository) UpdatePermission(ctx context.Context, id int, updates *user.PermissionUpdate) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if err := database.ExecuteUpdate(r.db, "ums.rbac_permission", updates, "id=$1", id); err != nil {
		return fmt.Errorf("updating permission id %d: %w", id, err)
	}
	return nil
}

func (r *repository) DeletePermission(ctx context.Context, id int) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if _, err := database.MustExec(r.db.ExecContext(ctx, "delete from ums.rbac_permission where id=$1", id)); err != nil {
		return fmt.Errorf("deleting permission id %d: %w", id, err)
	}
	return nil
}

func (r *repository) GetRolePermissions(ctx context.Context, id int) ([]rbac.Permission, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	permissions := []rbac.Permission{}
	if err := database.Retrieve(
		ctx,
		r.db,
		&permissions,
		`select b.id, b.code, b.description
		from ums.rbac_role_permission a
		join ums.rbac_permission b on a.permission_id=b.id 
		where a.role_id=$1
		order by b.id`,
		id); err != nil {
		return nil, fmt.Errorf("retrieving role permissions for role id %d: %w", id, err)
	}

	return permissions, nil
}

func (r *repository) AddRolePermission(ctx context.Context, roleID int, permissionID int) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if _, err := database.MustExec(r.db.ExecContext(ctx, `insert into ums.rbac_role_permission
		(role_id, permission_id)
		values
		($1, $2)`,
		roleID, permissionID)); err != nil {
		return fmt.Errorf("adding role permission for role id %d, permission id %d: %w", roleID, permissionID, err)
	}
	return nil
}

func (r *repository) DeleteRolePermission(ctx context.Context, roleID int, permissionID int) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if _, err := database.MustExec(r.db.ExecContext(ctx, `delete from ums.rbac_role_permission where role_id=$1 and permission_id=$2`, roleID, permissionID)); err != nil {
		return fmt.Errorf("deleting role permission for role id %d, permission id %d: %w", roleID, permissionID, err)
	}
	return nil
}
