package logic

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"

	"aqua/pkg/infra/logutil"
	blog "aqua/services/blog/domain"
)

func TestGetComments(t *testing.T) {
	tests := []struct {
		name        string
		aid         int
		page        int
		pageSize    int
		mockSetup   func(*commentRepoStub)
		expectError bool
	}{
		{
			name:     "Success",
			aid:      1,
			page:     1,
			pageSize: 10,
			mockSetup: func(s *commentRepoStub) {
				s.getCommentsFn = func(ctx context.Context, aid, offset, limit int) ([]blog.Comment, error) {
					return []blog.Comment{}, nil
				}
				s.getCommentCountFn = func(ctx context.Context, aid int) (int, error) { return 0, nil }
			},
			expectError: false,
		},
		{
			name:     "Repo Error",
			aid:      1,
			page:     1,
			pageSize: 10,
			mockSetup: func(s *commentRepoStub) {
				s.getCommentsFn = func(ctx context.Context, aid, offset, limit int) ([]blog.Comment, error) {
					return nil, errors.New("db error")
				}
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := &commentRepoStub{}
			if tc.mockSetup != nil {
				tc.mockSetup(repo)
			}
			svc := NewCommentService(repo)
			ctx := logutil.WithLogger(context.Background(), slog.Default())
			res, err := svc.GetComments(ctx, tc.aid, tc.page, tc.pageSize)
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
			}
		})
	}
}

func TestAddComment(t *testing.T) {
	tests := []struct {
		name        string
		aid         int
		uid         int
		comment     *blog.Comment
		mockSetup   func(*commentRepoStub)
		expectError bool
	}{
		{
			name:    "Success",
			aid:     1,
			uid:     123,
			comment: &blog.Comment{Content: "Nice post!"},
			mockSetup: func(s *commentRepoStub) {
				s.addCommentFn = func(ctx context.Context, c *blog.Comment) error {
					assert.Equal(t, 123, c.UID)
					assert.Equal(t, 1, c.Aid)
					return nil
				}
			},
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := &commentRepoStub{}
			if tc.mockSetup != nil {
				tc.mockSetup(repo)
			}
			svc := NewCommentService(repo)
			ctx := logutil.WithLogger(context.Background(), slog.Default())
			err := svc.AddComment(ctx, tc.aid, tc.uid, tc.comment)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
