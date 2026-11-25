package logic

import (
	"context"
	"encoding/json"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	"github.com/ipuppet/aqua/server/pkg/infra/logutil"
	auth "github.com/ipuppet/aqua/server/services/passport/domain"
)

func TestLoginByPassword(t *testing.T) {
	hashPwd := func(p string) string {
		b, _ := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
		return string(b)
	}

	// 1. Arrange: Setup the In-Memory Database and dependencies (Given)
	repo := newMemoryAuthRepo()
	validUserPwd := hashPwd("correctPassword123")

	repo.AddUser(&auth.Account{
		UID:      1,
		Username: "testuser",
		Email:    "test@example.com",
		Password: validUserPwd,
		Lock:     0,
	})

	repo.AddUser(&auth.Account{
		UID:      2,
		Username: "lockeduser",
		Password: validUserPwd,
		Lock:     1, // 1 implies locked
	})

	svc := NewAuthService(repo, &memoryRBACReader{}, mustTestPrivateKey(t))
	ctx := logutil.WithLogger(context.Background(), slog.Default())

	tests := []struct {
		name        string
		account     string
		password    string
		expectError bool
	}{
		{
			name:        "Successful Login with Username",
			account:     "testuser",
			password:    "correctPassword123",
			expectError: false,
		},
		{
			name:        "Successful Login with Email",
			account:     "test@example.com",
			password:    "correctPassword123",
			expectError: false,
		},
		{
			name:        "Invalid Account - Not Found",
			account:     "nonexistent",
			password:    "correctPassword123",
			expectError: true,
		},
		{
			name:        "Locked Account Disallowed",
			account:     "lockeduser",
			password:    "correctPassword123",
			expectError: true,
		},
		{
			name:        "Incorrect Password",
			account:     "testuser",
			password:    "wrongpass",
			expectError: true,
		},
		{
			name:        "Empty Password",
			account:     "testuser",
			password:    "",
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// 2. Act: Execute the authentication call (When)
			res, err := svc.LoginByPassword(ctx, tc.account, tc.password, LoginContext{
				LoginIP:    "203.0.113.8",
				DeviceInfo: "ios/safari",
			})

			// 3. Assert: Validate results (Then)
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.NotEmpty(t, res.AccessToken)
				assert.NotEmpty(t, res.RefreshToken)
				payload, marshalErr := json.Marshal(res)
				assert.NoError(t, marshalErr)
				assert.NotContains(t, string(payload), "public")
			}
		})
	}

	// Ensure the successful logins triggered the logging behavior in the DB
	t.Run("Verify Login Side Effects", func(t *testing.T) {
		deadline := time.Now().Add(500 * time.Millisecond)
		for time.Now().Before(deadline) {
			repo.RLock()
			auditCount := len(repo.loginAudits)
			repo.RUnlock()
			if auditCount > 0 {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}

		repo.RLock()
		defer repo.RUnlock()
		assert.NotEmpty(t, repo.loginAudits, "Expected successful logins to be recorded by RecordLoginSuccess")
		assert.Equal(t, "password", repo.loginAudits[0].LoginMethod)
		assert.Equal(t, "ios/safari", repo.loginAudits[0].DeviceInfo)
		assert.Equal(t, "203.0.113.8", repo.loginAudits[0].LoginIP)
	})
}
