package logic

import (
	"context"
	"io"
	"log/slog"
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ipuppet/aqua/server/pkg/infra/logutil"
	user "github.com/ipuppet/aqua/server/services/user/domain"
)

const (
	refreshTokenRevokeAfterUserKeyPrefix = "refresh_token_revoke_after:user:"
	refreshTokenRevokeAfterGlobalKey     = "refresh_token_revoke_after:global"
)

type setRecorderCache struct {
	setKeys []string
}

func (c *setRecorderCache) Get(ctx context.Context, key string) ([]byte, error) { return nil, nil }
func (c *setRecorderCache) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	c.setKeys = append(c.setKeys, key)
	return nil
}
func (c *setRecorderCache) Del(ctx context.Context, key string) error          { return nil }
func (c *setRecorderCache) DelMatch(ctx context.Context, pattern string) error { return nil }
func (c *setRecorderCache) GetMatch(ctx context.Context, pattern string) (map[string][]byte, error) {
	return make(map[string][]byte), nil
}
func (c *setRecorderCache) Flush(ctx context.Context) error { return nil }

func TestAddUserRole(t *testing.T) {
	tests := []struct {
		name        string
		uid         int
		id          int
		mockSetup   func(*managementRepoStub)
		expectError bool
	}{
		{
			name: "Success",
			uid:  1,
			id:   2,
			mockSetup: func(s *managementRepoStub) {
				s.addUserRoleFn = func(_ context.Context, uid, id int) error {
					assert.Equal(t, 1, uid)
					assert.Equal(t, 2, id)
					return nil
				}
			},
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := &managementRepoStub{}
			if tc.mockSetup != nil {
				tc.mockSetup(repo)
			}
			svc := NewAccessControlService(repo, nil)
			ctx := logutil.WithLogger(context.Background(), slog.New(slog.NewJSONHandler(io.Discard, nil)))
			err := svc.AddUserRole(ctx, tc.uid, tc.id)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdateRole_TargetedRefreshRevocation(t *testing.T) {
	recorder := &setRecorderCache{}
	repo := &managementRepoStub{}
	repo.cache = recorder
	repo.getUserIDsByRoleFn = func(_ context.Context, id int) ([]int, error) {
		assert.Equal(t, 9, id)
		return []int{3, 7}, nil
	}
	repo.updateRoleFn = func(_ context.Context, id int, updates *user.RoleUpdate) error {
		assert.Equal(t, 9, id)
		assert.NotNil(t, updates)
		return nil
	}

	svc := NewAccessControlService(repo, nil)
	ctx := logutil.WithLogger(context.Background(), slog.New(slog.NewJSONHandler(io.Discard, nil)))
	code := "auditor"
	err := svc.UpdateRole(ctx, 9, &user.RoleUpdate{Code: &code})
	assert.NoError(t, err)

	sort.Strings(recorder.setKeys)
	assert.Equal(t, []string{
		refreshTokenRevokeAfterUserKeyPrefix + strconv.Itoa(3),
		refreshTokenRevokeAfterUserKeyPrefix + strconv.Itoa(7),
	}, recorder.setKeys)
	assert.NotContains(t, recorder.setKeys, refreshTokenRevokeAfterGlobalKey)
}
