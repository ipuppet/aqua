package dao

import (
	"context"
	"fmt"
	"time"

	auth "aqua/services/passport/domain"
)

// RecordLoginLog updates user last_login_at and inserts an audit record in one transaction.
func (r *repository) RecordLoginLog(ctx context.Context, record auth.LoginAuditRecord) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if record.CreatedAt.IsZero() {
		record.CreatedAt = time.Now().UTC()
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx,
		"update ums.rbac_user set last_login_at=$1 where uid=$2",
		record.CreatedAt, record.UID,
	); err != nil {
		_ = tx.Rollback()
		return err
	}

	if _, err := tx.ExecContext(ctx,
		`insert into ums.user_login_audit (uid, login_ip, device_info, login_method, refresh_jti, created_at)
		 values ($1, $2, $3, $4, $5, $6)`,
		record.UID,
		record.LoginIP,
		record.DeviceInfo,
		record.LoginMethod,
		record.RefreshJTI,
		record.CreatedAt,
	); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing login audit transaction: %w", err)
	}

	// Keep user cache and table in sync for last_login_at after successful transaction.
	_ = r.updateCachedLastLoginAt(ctx, record.UID, record.CreatedAt)
	return nil
}

// ListRecentLogins retrieves recent login audit records for a user.
func (r *repository) ListRecentLogins(ctx context.Context, uid int, limit int) ([]auth.LoginAuditRecord, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if limit <= 0 {
		limit = 100
	}

	rows, err := r.db.QueryContext(ctx,
		`select uid, login_ip, device_info, login_method, refresh_jti, created_at
		 from ums.user_login_audit
		 where uid = $1
		 order by created_at desc
		 limit $2`,
		uid, limit,
	)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	out := make([]auth.LoginAuditRecord, 0)
	for rows.Next() {
		var record auth.LoginAuditRecord
		if err := rows.Scan(
			&record.UID,
			&record.LoginIP,
			&record.DeviceInfo,
			&record.LoginMethod,
			&record.RefreshJTI,
			&record.CreatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, record)
	}
	return out, rows.Err()
}
