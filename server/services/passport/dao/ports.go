package dao

import (
	"context"
	"time"

	"github.com/ipuppet/aqua/server/pkg/infra/cache"
	"github.com/ipuppet/aqua/server/pkg/infra/database"
	auth "github.com/ipuppet/aqua/server/services/passport/domain"
)

// AuthRepository exposes all persistence operations used by the passport logic layer.
type AuthRepository interface {
	Cache() cache.Store
	GetUserByAccount(ctx context.Context, account string) (*auth.Account, error)
	GetUserByUid(ctx context.Context, uid int) (*auth.Account, error)
	RecordLoginLog(ctx context.Context, record auth.LoginAuditRecord) error
	RegisterUser(ctx context.Context, um *auth.Account) error
	StoreVerificationCode(ctx context.Context, key, code string) error
	DeleteVerificationCode(ctx context.Context, key string) error
	GetVerificationCode(ctx context.Context, key string) (string, error)
	StoreUserSession(ctx context.Context, session *auth.Session, ttl time.Duration) error
	GetUserSession(ctx context.Context, uid int, sessionID string) (*auth.Session, error)
	GetUserSessions(ctx context.Context, uid int) ([]auth.Session, error)
	DeleteUserSession(ctx context.Context, uid int, sessionID string) error
	DeleteAllUserSessions(ctx context.Context, uid int) error
	ListRecentLogins(ctx context.Context, uid int, limit int) ([]auth.LoginAuditRecord, error)
	UpdateUserPassword(ctx context.Context, uid int, password string) error
	StorePasswordResetToken(ctx context.Context, token, email string) error
	GetPasswordResetEmail(ctx context.Context, token string) (string, error)
	DeletePasswordResetToken(ctx context.Context, token string) error
}

// repository is the concrete implementation backed by PostgreSQL and Redis.
type repository struct {
	db    *database.Client
	cache cache.Store
}

// NewRepository builds a new repository using the provided dependencies.
func NewRepository(db *database.Client, c cache.Store) *repository {
	return &repository{db: db, cache: c}
}

func (r *repository) Cache() cache.Store {
	return r.cache
}

var _ AuthRepository = (*repository)(nil)
