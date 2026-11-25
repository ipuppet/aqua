package logic

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"aqua/pkg/infra/logutil"
	blog "aqua/services/blog/domain"
)

func TestGetArticleDetail(t *testing.T) {
	tests := []struct {
		name        string
		aid         int
		ip          string
		mockSetup   func(*articleRepoStub, *accessRepoStub)
		expectError bool
	}{
		{
			name: "Success",
			aid:  1,
			ip:   "127.0.0.1",
			mockSetup: func(articleRepo *articleRepoStub, accessRepo *accessRepoStub) {
				articleRepo.getLastAidFn = func(ctx context.Context) (int, error) { return 10, nil }
				articleRepo.getContextFn = func(ctx context.Context, startAid, direction, boundary int) (*blog.ArticleContext, error) {
					return &blog.ArticleContext{Aid: 2, Title: "Next"}, nil
				}
				articleRepo.getArticleFn = func(ctx context.Context, aid int) (*blog.Article, error) {
					return &blog.Article{Aid: 1, Title: "Current"}, nil
				}
				accessRepo.logAccessInfoFn = func(ctx context.Context, aid int, ip string) error { return nil }
			},
			expectError: false,
		},
		{
			name: "GetLastAid Error",
			aid:  1,
			ip:   "127.0.0.1",
			mockSetup: func(articleRepo *articleRepoStub, _ *accessRepoStub) {
				articleRepo.getLastAidFn = func(ctx context.Context) (int, error) { return 0, errors.New("db error") }
			},
			expectError: true,
		},
		{
			name: "GetArticle Error",
			aid:  1,
			ip:   "127.0.0.1",
			mockSetup: func(articleRepo *articleRepoStub, _ *accessRepoStub) {
				articleRepo.getLastAidFn = func(ctx context.Context) (int, error) { return 10, nil }
				articleRepo.getArticleFn = func(ctx context.Context, aid int) (*blog.Article, error) {
					return nil, errors.New("not found")
				}
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := &articleRepoStub{}
			accessRepo := &accessRepoStub{}
			if tc.mockSetup != nil {
				tc.mockSetup(repo, accessRepo)
			}
			svc := NewArticleService(repo, accessRepo, &indexerStub{})
			ctx := logutil.WithLogger(context.Background(), slog.Default())
			res, err := svc.GetArticleDetail(ctx, tc.aid, tc.ip)
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, "Current", res.Article.Title)
			}
		})
	}
}

func TestGetArticleDetailLogsAccessAfterRequestCancel(t *testing.T) {
	logged := make(chan struct {
		aid int
		ip  string
	}, 1)
	hasDeadline := make(chan bool, 1)
	repo := &articleRepoStub{
		getLastAidFn: func(ctx context.Context) (int, error) { return 10, nil },
		getArticleFn: func(ctx context.Context, aid int) (*blog.Article, error) {
			return &blog.Article{Aid: 1, Title: "Current"}, nil
		},
	}
	accessRepo := &accessRepoStub{
		logAccessInfoFn: func(ctx context.Context, aid int, ip string) error {
			_, ok := ctx.Deadline()
			hasDeadline <- ok
			logged <- struct {
				aid int
				ip  string
			}{aid: aid, ip: ip}
			return nil
		},
	}
	svc := NewArticleService(repo, accessRepo, &indexerStub{})
	ctx, cancel := context.WithCancel(logutil.WithLogger(context.Background(), slog.Default()))
	cancel()

	res, err := svc.GetArticleDetail(ctx, 1, "127.0.0.1")
	require.NoError(t, err)
	require.NotNil(t, res)

	select {
	case got := <-logged:
		assert.Equal(t, 1, got.aid)
		assert.Equal(t, "127.0.0.1", got.ip)
		assert.True(t, <-hasDeadline)
	case <-time.After(time.Second):
		t.Fatal("access log job was not called")
	}
}

func TestGetArticles(t *testing.T) {
	tests := []struct {
		name        string
		page        int
		pageSize    int
		mockSetup   func(*articleRepoStub)
		expectError bool
	}{
		{
			name:     "Pagination Defaults",
			page:     0,
			pageSize: 0,
			mockSetup: func(s *articleRepoStub) {
				s.getArticlesFn = func(ctx context.Context, offset, limit int) ([]blog.Article, error) {
					assert.Equal(t, 0, offset)
					assert.Equal(t, 10, limit)
					return []blog.Article{}, nil
				}
				s.getCountFn = func(ctx context.Context) (int, error) { return 0, nil }
			},
			expectError: false,
		},
		{
			name:     "Custom Pagination",
			page:     2,
			pageSize: 5,
			mockSetup: func(s *articleRepoStub) {
				s.getArticlesFn = func(ctx context.Context, offset, limit int) ([]blog.Article, error) {
					assert.Equal(t, 5, offset)
					assert.Equal(t, 5, limit)
					return []blog.Article{}, nil
				}
				s.getCountFn = func(ctx context.Context) (int, error) { return 0, nil }
			},
			expectError: false,
		},
		{
			name:     "Repo Error",
			page:     1,
			pageSize: 10,
			mockSetup: func(s *articleRepoStub) {
				s.getArticlesFn = func(ctx context.Context, offset, limit int) ([]blog.Article, error) {
					return nil, errors.New("db error")
				}
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := &articleRepoStub{}
			accessRepo := &accessRepoStub{}
			if tc.mockSetup != nil {
				tc.mockSetup(repo)
			}
			svc := NewArticleService(repo, accessRepo, &indexerStub{})
			res, err := svc.GetArticles(context.Background(), tc.page, tc.pageSize)
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

func TestSearchArticles(t *testing.T) {
	tests := []struct {
		name        string
		kw          string
		page        int
		pageSize    int
		indexer     *indexerStub
		expectError bool
	}{
		{
			name: "Search Success",
			kw:   "Go",
			page: 1, pageSize: 10,
			indexer: &indexerStub{
				searchArticlesFn: func(_ context.Context, query string, limit, offset int) (*blog.SearchResponse, error) {
					assert.Equal(t, "Go", query)
					assert.Equal(t, 10, limit)
					assert.Equal(t, 0, offset)
					return &blog.SearchResponse{
						Items: []blog.Article{{Aid: 1, Title: "Go Tutorial"}},
						Total: 1,
					}, nil
				},
			},
			expectError: false,
		},
		{
			name: "Search Indexer Error",
			kw:   "foo", page: 1, pageSize: 5,
			indexer: &indexerStub{
				searchArticlesFn: func(context.Context, string, int, int) (*blog.SearchResponse, error) {
					return nil, errors.New("meili error")
				},
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := &articleRepoStub{}
			accessRepo := &accessRepoStub{}
			idx := tc.indexer
			if idx == nil {
				idx = &indexerStub{}
			}
			svc := NewArticleService(repo, accessRepo, idx)
			res, err := svc.SearchArticles(context.Background(), tc.kw, tc.page, tc.pageSize)
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, 1, res.Total)
				assert.Len(t, res.Items, 1)
				assert.Equal(t, "Go Tutorial", res.Items[0].Title)
			}
		})
	}
}
