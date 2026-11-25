package rbac

import (
	"context"
	"database/sql"
	"time"

	"github.com/ipuppet/aqua/server/pkg/infra/cache"
	"github.com/ipuppet/aqua/server/pkg/infra/database"
)

const (
	userRolesTTL = 24 * time.Hour
	userPermsTTL = 24 * time.Hour
)

// Reader provides canonical RBAC read operations.
type Reader interface {
	LoadUserRoles(ctx context.Context, uid int) ([]Role, error)
	LoadUserPermissions(ctx context.Context, uid int) ([]Permission, error)
	LoadRolePermissions(ctx context.Context, role string) ([]Permission, error)
}

type reader struct {
	db    *database.Client
	cache cache.Store
}

// NewReader builds the RBAC reader.
func NewReader(db *database.Client, c cache.Store) Reader {
	return &reader{db: db, cache: c}
}

func (r *reader) LoadUserRoles(ctx context.Context, uid int) ([]Role, error) {
	var roles []Role
	err := database.FetchWithCache(ctx, r.cache, CacheKeyUserRoles(uid), userRolesTTL, &roles,
		func() (any, error) {
			list := []Role{}
			err := database.Retrieve(ctx, r.db, &list,
				`SELECT b.code, b.description, b.id
				FROM ums.rbac_role b
				JOIN ums.rbac_user_role c ON b.id = c.role_id
				WHERE c.uid = $1`, uid)
			return list, err
		})
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if roles == nil {
		roles = []Role{}
	}
	return roles, nil
}

func (r *reader) LoadUserPermissions(ctx context.Context, uid int) ([]Permission, error) {
	roles, err := r.LoadUserRoles(ctx, uid)
	if err != nil {
		return nil, err
	}

	seen := make(map[string]struct{})
	perms := make([]Permission, 0, len(roles))
	for _, role := range roles {
		rolePerms, err := r.LoadRolePermissions(ctx, role.Code)
		if err != nil {
			return nil, err
		}
		for _, p := range rolePerms {
			if _, ok := seen[p.Code]; !ok {
				seen[p.Code] = struct{}{}
				perms = append(perms, p)
			}
		}
	}
	return perms, nil
}

func (r *reader) LoadRolePermissions(ctx context.Context, role string) ([]Permission, error) {
	var perms []Permission
	err := database.FetchWithCache(ctx, r.cache, CacheKeyStringRolePermissions(role), userPermsTTL, &perms,
		func() (any, error) {
			list := []Permission{}
			err := database.Retrieve(ctx, r.db, &list,
				`SELECT b.id, b.code, b.description
				FROM ums.rbac_role_permission a
				JOIN ums.rbac_permission b ON a.permission_id = b.id
				JOIN ums.rbac_role c ON a.role_id = c.id
					WHERE c.code = $1`, role)
			return list, err
		})
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if perms == nil {
		perms = []Permission{}
	}
	return perms, nil
}
