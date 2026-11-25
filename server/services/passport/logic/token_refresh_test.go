package logic

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/ipuppet/aqua/server/pkg/transport/token"
	authdomain "github.com/ipuppet/aqua/server/services/passport/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"github.com/ipuppet/aqua/server/pkg/infra/logutil"
)

func TestVerifyAndRefreshToken_DoesNotWriteLoginAudit(t *testing.T) {
	hashed, err := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.DefaultCost)
	require.NoError(t, err)

	repo := newMemoryAuthRepo()
	repo.AddUser(&authdomain.Account{
		UID:      18,
		Username: "refresh_user",
		Email:    "refresh@example.com",
		Password: string(hashed),
		Lock:     0,
	})

	svc := NewAuthService(repo, &memoryRBACReader{}, mustTestPrivateKey(t))
	ctx := logutil.WithLogger(context.Background(), slog.Default())

	loginRes, err := svc.LoginByPassword(ctx, "refresh_user", "secret123", LoginContext{
		LoginIP:    "198.51.100.12",
		DeviceInfo: "android/chrome",
	})
	require.NoError(t, err)
	oldRefreshToken := loginRes.RefreshToken
	require.NotEmpty(t, oldRefreshToken)

	publicKey := &svc.privateKey.PublicKey
	oldClaims, err := token.VerifyToken(publicKey, &token.RefreshTokenClaims{}, oldRefreshToken)
	require.NoError(t, err)

	deadline := time.Now().Add(500 * time.Millisecond)
	for time.Now().Before(deadline) {
		repo.RLock()
		count := len(repo.loginAudits)
		repo.RUnlock()
		if count > 0 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	repo.RLock()
	preRefreshAuditCount := len(repo.loginAudits)
	repo.RUnlock()
	require.Equal(t, 1, preRefreshAuditCount)

	refreshed, err := svc.VerifyAndRefreshToken(ctx, "", oldRefreshToken)
	require.NoError(t, err)
	newRefreshToken := refreshed.RefreshToken
	require.NotEmpty(t, newRefreshToken)

	newClaims, err := token.VerifyToken(publicKey, &token.RefreshTokenClaims{}, newRefreshToken)
	require.NoError(t, err)
	assert.Equal(t, oldClaims.GetSessionID(), newClaims.GetSessionID())
	assert.NotEqual(t, oldClaims.TokenID(), newClaims.TokenID())

	oldSession, err := repo.GetUserSession(ctx, int(oldClaims.UID), oldClaims.GetSessionID())
	require.NoError(t, err)
	assert.Equal(t, "android/chrome", oldSession.DeviceInfo)

	newSession, err := repo.GetUserSession(ctx, int(newClaims.UID), newClaims.GetSessionID())
	require.NoError(t, err)
	assert.Equal(t, oldSession.UID, newSession.UID)
	assert.Equal(t, oldSession.SessionID, newSession.SessionID)
	assert.Equal(t, "android/chrome", newSession.DeviceInfo)
	assert.False(t, newSession.CreatedAt.IsZero())

	repo.RLock()
	postRefreshAuditCount := len(repo.loginAudits)
	repo.RUnlock()
	assert.Equal(t, preRefreshAuditCount, postRefreshAuditCount)
}

func TestVerifyAndRefreshToken_WithValidAccessToken_DoesNotReturnPublic(t *testing.T) {
	hashed, err := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.DefaultCost)
	require.NoError(t, err)

	repo := newMemoryAuthRepo()
	repo.AddUser(&authdomain.Account{
		UID:      21,
		Username: "access_user",
		Email:    "access@example.com",
		Password: string(hashed),
		Lock:     0,
	})

	svc := NewAuthService(repo, &memoryRBACReader{}, mustTestPrivateKey(t))
	ctx := logutil.WithLogger(context.Background(), slog.Default())

	loginRes, err := svc.LoginByPassword(ctx, "access_user", "secret123", LoginContext{
		LoginIP:    "203.0.113.10",
		DeviceInfo: "ios/safari",
	})
	require.NoError(t, err)
	accessToken := loginRes.AccessToken
	require.NotEmpty(t, accessToken)

	res, err := svc.VerifyAndRefreshToken(ctx, accessToken, "")
	require.NoError(t, err)
	assert.Empty(t, res.AccessToken)
	assert.Empty(t, res.RefreshToken)
}
