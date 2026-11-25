package logic

import (
	"context"
	"crypto/rsa"
	"log/slog"
	"time"

	"aqua/pkg/infra/logutil"
	"aqua/pkg/rbac"
	"aqua/services/passport/dao"
)

const (
	JwtIssuer        = "aqua"
	accessTokenLife  = 15 * time.Minute
	refreshTokenLife = 7 * 24 * time.Hour
	// refreshTokenGracePeriod is the TTL applied to a rotated refresh token so that
	// concurrent in-flight requests using the same old token do not get unexpectedly
	// rejected during the brief window between rotation and client acknowledgement.
	refreshTokenGracePeriod = 15 * time.Second
)

// LoginContext carries request metadata captured at login time.
type LoginContext struct {
	LoginIP    string
	DeviceInfo string
}

// LoginResult is the response payload for password login.
type LoginResult struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Exp          int64  `json:"exp"`
}

// TokenRefreshResult is the response payload for token refresh.
// Empty fields mean the current access token is still valid and no rotation happened.
type TokenRefreshResult struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// AuthService is the default AuthService implementation.
type AuthService struct {
	repo       dao.AuthRepository
	rbacReader rbac.Reader
	privateKey *rsa.PrivateKey
}

// NewAuthService wires dependencies into an AuthService implementation.
func NewAuthService(repo dao.AuthRepository, rbacReader rbac.Reader, privateKey *rsa.PrivateKey) *AuthService {
	return &AuthService{repo: repo, rbacReader: rbacReader, privateKey: privateKey}
}

func (s *AuthService) log(ctx context.Context) *slog.Logger {
	return logutil.FromContext(ctx)
}
