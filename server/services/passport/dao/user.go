package dao

import (
	"context"
	"database/sql"
	"encoding/json"
	"regexp"
	"time"

	"github.com/ipuppet/aqua/server/pkg/infra/database"
	"github.com/ipuppet/aqua/server/pkg/rbac"
	auth "github.com/ipuppet/aqua/server/services/passport/domain"
	"github.com/ipuppet/aqua/server/services/passport/internal/regex"
)

const userInfoTTL = 24 * time.Hour

// GetUserByAccount fetches a user model by account identifier which can be username, email, or phone.
// Bypasses cache to guarantee fresh credentials for authentication.
func (r *repository) GetUserByAccount(ctx context.Context, account string) (*auth.Account, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	var column string
	accountByte := []byte(account)
	if matched, _ := regexp.Match(regex.Email, accountByte); matched {
		column = "email"
	} else if matched, _ := regexp.Match(regex.PhoneCN, accountByte); matched {
		column = "phone"
	} else {
		column = "username"
	}

	model := &auth.Account{}
	err := database.Retrieve(ctx, r.db, model,
		`SELECT `+auth.AccountColumns+` FROM ums.rbac_user WHERE `+column+`=$1 LIMIT 1`, account)
	if err != nil {
		return nil, err
	}
	return model, nil
}

// GetUserByUid fetches a user model via uid with caching.
func (r *repository) GetUserByUid(ctx context.Context, uid int) (*auth.Account, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	model := &auth.Account{}
	err := database.FetchWithCache(ctx, r.cache, rbac.CacheKeyUserInfo(uid), userInfoTTL, model,
		func() (any, error) {
			m := &auth.Account{}
			err := database.Retrieve(ctx, r.db, m,
				`SELECT `+auth.AccountColumns+` FROM ums.rbac_user WHERE uid=$1 LIMIT 1`, uid)
			return m, err
		})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}
	return model, nil
}

func (r *repository) updateCachedLastLoginAt(ctx context.Context, uid int, lastLoginAt time.Time) error {
	if r.cache == nil {
		return nil
	}
	key := rbac.CacheKeyUserInfo(uid)
	b, err := r.cache.Get(ctx, key)
	if err != nil {
		return err
	}
	var model auth.Account
	if err := json.Unmarshal(b, &model); err != nil {
		return err
	}
	model.LastLoginAt = &lastLoginAt
	return r.cache.Set(ctx, key, model, userInfoTTL)
}

// RegisterUser inserts a new user and grants default role membership inside a single transaction.
func (r *repository) RegisterUser(ctx context.Context, um *auth.Account) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	var uid int
	if err := tx.QueryRowContext(
		ctx,
		"insert into ums.rbac_user (username, nickname, email, password, gender) values ($1, $2, $3, $4, $5) RETURNING uid",
		um.Username, um.Nickname, um.Email, um.Password, um.Gender,
	).Scan(&uid); err != nil {
		_ = tx.Rollback()
		return err
	}

	if _, err := database.MustExec(tx.ExecContext(
		ctx,
		"insert into ums.rbac_user_role (uid, role_id) values ($1, $2)",
		uid, 0,
	)); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

// UpdateUserPassword updates the password of the user identified by uid.
func (r *repository) UpdateUserPassword(ctx context.Context, uid int, password string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	_, err := r.db.ExecContext(ctx, "update ums.rbac_user set password=$1 where uid=$2", password, uid)
	return err
}
