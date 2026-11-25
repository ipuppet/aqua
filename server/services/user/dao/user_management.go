package dao

import (
	"context"
	"fmt"
	"strings"

	"aqua/pkg/infra/database"
	user "aqua/services/user/domain"
)

func (r *repository) SearchUsers(ctx context.Context, kw string, offset int, limit int) ([]user.UserModel, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var users []user.UserModel
	err := database.Retrieve(
		ctx,
		r.db,
		&users,
		`select
		uid, username, nickname, email, phone, avatar, gender, lock, register_at, last_login_at
		from ums.rbac_user
		where username like $1 escape '\\'
		limit $2 offset $3`,
		"%"+escapeLikePattern(kw)+"%", limit, offset)
	if err != nil {
		return nil, fmt.Errorf("searching users with keyword '%s', offset %d, limit %d: %w", kw, offset, limit, err)
	}

	return users, nil
}

// escapeLikePattern escapes SQL LIKE special characters in s so that they are
// treated as literals rather than wildcards. The caller must pair this with
// ESCAPE '\\' in the SQL query.
func escapeLikePattern(s string) string {
	s = strings.ReplaceAll(s, `\\`, `\\\\`)
	s = strings.ReplaceAll(s, `%`, `\\%`)
	s = strings.ReplaceAll(s, `_`, `\\_`)
	return s
}

func (r *repository) GetUsers(ctx context.Context, offset int, limit int) ([]user.UserModel, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var users []user.UserModel
	err := database.Retrieve(
		ctx,
		r.db,
		&users,
		`select
		uid, username, nickname, email, phone, avatar, gender, lock, register_at, last_login_at
		from ums.rbac_user
		limit $1 offset $2`,
		limit, offset)
	if err != nil {
		return nil, fmt.Errorf("retrieving users with offset %d, limit %d: %w", offset, limit, err)
	}

	return users, nil
}
