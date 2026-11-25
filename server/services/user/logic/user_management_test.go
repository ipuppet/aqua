package logic

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"

	"aqua/pkg/infra/logutil"
	user "aqua/services/user/domain"
)

func TestGetUsers(t *testing.T) {
	tests := []struct {
		name           string
		page           int
		pageSize       int
		expectedOffset int
		expectedLimit  int
		mockSetup      func(*managementRepoStub)
		expectError    bool
	}{
		{
			name:           "Pagination Normalization",
			page:           0,
			pageSize:       0,
			expectedOffset: 0,
			expectedLimit:  10,
			mockSetup: func(s *managementRepoStub) {
				s.getUsersFn = func(_ context.Context, offset, limit int) ([]user.UserModel, error) {
					assert.Equal(t, 0, offset)
					assert.Equal(t, 10, limit)
					return []user.UserModel{{UID: 1}}, nil
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
			svc := NewUserManagementService(repo)
			ctx := logutil.WithLogger(context.Background(), slog.New(slog.NewJSONHandler(io.Discard, nil)))
			res, err := svc.GetUsers(ctx, tc.page, tc.pageSize)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, res)
			}
		})
	}
}

func TestSearchUsers(t *testing.T) {
	tests := []struct {
		name        string
		kw          string
		mockSetup   func(*managementRepoStub)
		expectError bool
	}{
		{
			name: "Success",
			kw:   "test",
			mockSetup: func(s *managementRepoStub) {
				s.searchUsersFn = func(_ context.Context, kw string, offset, limit int) ([]user.UserModel, error) {
					assert.Equal(t, "test", kw)
					return []user.UserModel{}, nil
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
			svc := NewUserManagementService(repo)
			ctx := logutil.WithLogger(context.Background(), slog.New(slog.NewJSONHandler(io.Discard, nil)))
			_, err := svc.SearchUsers(ctx, tc.kw, 1, 10)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	tests := []struct {
		name        string
		uid         int
		mockSetup   func(*managementRepoStub)
		expectError bool
	}{
		{
			name: "Success",
			uid:  1,
			mockSetup: func(s *managementRepoStub) {
				s.updateFn = func(_ context.Context, uid int, _ *user.UserUpdate) error {
					assert.Equal(t, 1, uid)
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
			svc := NewUserManagementService(repo)
			ctx := logutil.WithLogger(context.Background(), slog.New(slog.NewJSONHandler(io.Discard, nil)))
			err := svc.UpdateUser(ctx, tc.uid, &user.UserUpdate{})
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
