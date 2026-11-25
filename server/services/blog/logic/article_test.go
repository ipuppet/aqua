package logic

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"aqua/pkg/infra/logutil"
	blog "aqua/services/blog/domain"
)

// --- Tests ---

func TestCreateArticle(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		content     string
		tags        []blog.Tag
		mockSetup   func(*managementRepoStub)
		expectError bool
	}{
		{
			name:    "CA-01: Base Case - Create with Tags",
			title:   "My First Post",
			content: "# Hello World",
			tags:    []blog.Tag{{Name: "Go"}, {Name: "Testing"}},
			mockSetup: func(s *managementRepoStub) {
				s.createArticleFn = func(ctx context.Context, a *blog.Article) (int, error) {
					return 1, nil
				}
			},
			expectError: false,
		},
		{
			name:    "CA-02: Create without Tags",
			title:   "My First Post",
			content: "# Hello World",
			tags:    []blog.Tag{},
			mockSetup: func(s *managementRepoStub) {
				s.createArticleFn = func(ctx context.Context, a *blog.Article) (int, error) {
					return 2, nil
				}
			},
			expectError: false,
		},
		{
			name:    "CA-03: Empty Title (Simulated DB Error)",
			title:   "",
			content: "# Hello World",
			tags:    []blog.Tag{{Name: "Go"}},
			mockSetup: func(s *managementRepoStub) {
				s.createArticleFn = func(ctx context.Context, a *blog.Article) (int, error) {
					if a.Title == "" {
						return 0, errors.New("title cannot be empty")
					}
					return 3, nil
				}
			},
			expectError: true,
		},
		{
			name:    "CA-04: Empty Content",
			title:   "Empty Content Post",
			content: "",
			tags:    []blog.Tag{{Name: "Go"}},
			mockSetup: func(s *managementRepoStub) {
				s.createArticleFn = func(ctx context.Context, a *blog.Article) (int, error) {
					return 4, nil
				}
			},
			expectError: false,
		},
		{
			name:    "CA-05: Very Long Title (Simulated DB Error)",
			title:   strings.Repeat("A", 1000), // Assuming DB limit is < 1000
			content: "Content",
			tags:    []blog.Tag{},
			mockSetup: func(s *managementRepoStub) {
				s.createArticleFn = func(ctx context.Context, a *blog.Article) (int, error) {
					if len(a.Title) > 255 {
						return 0, errors.New("title too long")
					}
					return 5, nil
				}
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			adminRepo := &managementRepoStub{}
			articleRepo := &articleRepoStub{}
			if tc.mockSetup != nil {
				tc.mockSetup(adminRepo)
			}
			svc := NewManagementService(adminRepo, adminRepo, articleRepo, &indexerStub{})
			ctx := logutil.WithLogger(context.Background(), slog.Default())
			article := &blog.Article{
				Title:    tc.title,
				Markdown: tc.content,
				Tags:     tc.tags,
			}
			aid, err := svc.CreateArticle(ctx, article)
			if tc.expectError {
				assert.Error(t, err)
				assert.Equal(t, 0, aid)
			} else {
				assert.NoError(t, err)
				assert.NotZero(t, aid)
			}
		})
	}
}

func TestCreateArticleIndexesAfterRequestCancel(t *testing.T) {
	indexed := make(chan *blog.Article, 1)
	hasDeadline := make(chan bool, 1)
	idx := &indexerStub{
		indexArticleFn: func(ctx context.Context, article *blog.Article) error {
			_, ok := ctx.Deadline()
			hasDeadline <- ok
			indexed <- article
			return nil
		},
	}
	adminRepo := &managementRepoStub{
		createArticleFn: func(ctx context.Context, a *blog.Article) (int, error) {
			return 17, nil
		},
	}
	svc := NewManagementService(adminRepo, adminRepo, &articleRepoStub{}, idx)
	ctx, cancel := context.WithCancel(logutil.WithLogger(context.Background(), slog.Default()))
	cancel()

	article := &blog.Article{Title: "Indexed", Markdown: "body"}
	aid, err := svc.CreateArticle(ctx, article)
	require.NoError(t, err)
	require.Equal(t, 17, aid)

	select {
	case got := <-indexed:
		assert.Equal(t, 17, got.Aid)
		assert.Equal(t, "Indexed", got.Title)
		assert.True(t, <-hasDeadline)
	case <-time.After(time.Second):
		t.Fatal("search index job was not called")
	}
}
