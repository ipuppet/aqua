package dao

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"aqua/pkg/infra/database"
	"aqua/pkg/rbac"
	user "aqua/services/user/domain"
	schemauser "aqua/services/user/domain/schema"
)

const userInfoTTL = 24 * time.Hour

func (r *repository) GetUserInfo(ctx context.Context, uid int) (*user.UserModel, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	model := &user.UserModel{}
	err := database.FetchWithCache(ctx, r.cache, rbac.CacheKeyUserInfo(uid), userInfoTTL, model,
		func() (any, error) {
			m := &user.UserModel{}
			err := database.Retrieve(ctx, r.db, m,
				`SELECT `+schemauser.UserColumns+` FROM ums.rbac_user WHERE uid=$1 LIMIT 1`, uid)
			return m, err
		})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("getting user info for uid %d: %w", uid, err)
	}
	return model, nil
}

func (r *repository) UpdateUser(ctx context.Context, uid int, updates *user.UserUpdate) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if err := database.ExecuteUpdate(r.db, "ums.rbac_user", updates, "uid=$1", uid); err != nil {
		return fmt.Errorf("updating user %d: %w", uid, err)
	}
	return nil
}
