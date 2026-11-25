package logic

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"

	"aqua/pkg/infra/logutil"
	user "aqua/services/user/domain"
)

func TestGetUserInfo(t *testing.T) {
	tests := []struct {
		name        string
		uid         int
		mockSetup   func(*userRepoStub)
		expectError bool
	}{
		{
			name: "Success",
			uid:  42,
			mockSetup: func(s *userRepoStub) {
				s.getInfoFn = func(_ context.Context, uid int) (*user.UserModel, error) {
					if uid != 42 {
						return nil, errors.New("unexpected uid")
					}
					return &user.UserModel{UID: uid}, nil
				}
			},
			expectError: false,
		},
		{
			name: "Repo Error",
			uid:  42,
			mockSetup: func(s *userRepoStub) {
				s.getInfoFn = func(_ context.Context, uid int) (*user.UserModel, error) {
					return nil, errors.New("db error")
				}
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := &userRepoStub{}
			if tc.mockSetup != nil {
				tc.mockSetup(repo)
			}
			svc := NewUserService(repo, nil, nil)
			ctx := logutil.WithLogger(context.Background(), slog.New(slog.NewJSONHandler(io.Discard, nil)))
			resp, err := svc.GetUserInfo(ctx, tc.uid)
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tc.uid, resp.UID)
			}
		})
	}
}

func TestUpdateSelf(t *testing.T) {
	tests := []struct {
		name        string
		uid         int
		updates     *user.UserUpdate
		mockSetup   func(*userRepoStub)
		expectError bool
	}{
		{
			name:        "Missing User Context",
			uid:         0,
			updates:     &user.UserUpdate{},
			expectError: true,
		},
		{
			name:    "Success",
			uid:     7,
			updates: &user.UserUpdate{},
			mockSetup: func(s *userRepoStub) {
				s.updateFn = func(_ context.Context, uid int, _ *user.UserUpdate) error {
					assert.Equal(t, 7, uid)
					return nil
				}
			},
			expectError: false,
		},
		{
			name:    "Repo Error",
			uid:     7,
			updates: &user.UserUpdate{},
			mockSetup: func(s *userRepoStub) {
				s.updateFn = func(_ context.Context, uid int, _ *user.UserUpdate) error {
					return errors.New("db error")
				}
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := &userRepoStub{}
			if tc.mockSetup != nil {
				tc.mockSetup(repo)
			}
			svc := NewUserService(repo, nil, nil)
			ctx := logutil.WithLogger(context.Background(), slog.New(slog.NewJSONHandler(io.Discard, nil)))
			err := svc.UpdateSelf(ctx, tc.uid, tc.updates)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
