package logic

import (
	"context"
	"time"

	"aqua/pkg/infra/cache"
	"aqua/services/blog/dao"
	blog "aqua/services/blog/domain"
	"aqua/services/blog/meili"
)

// indexerStub implements meili.BlogSearchIndexer for tests.
type indexerStub struct {
	searchArticlesFn  func(context.Context, string, int, int) (*blog.SearchResponse, error)
	indexArticleFn    func(context.Context, *blog.Article) error
	deleteArticleFn   func(context.Context, int) error
	ensureBlogIndexFn func() error
}

func (s *indexerStub) SearchArticles(ctx context.Context, query string, limit, offset int) (*blog.SearchResponse, error) {
	if s.searchArticlesFn != nil {
		return s.searchArticlesFn(ctx, query, limit, offset)
	}
	return &blog.SearchResponse{Items: []blog.Article{}, Total: 0}, nil
}

func (s *indexerStub) IndexArticle(ctx context.Context, article *blog.Article) error {
	if s.indexArticleFn != nil {
		return s.indexArticleFn(ctx, article)
	}
	return nil
}

func (s *indexerStub) DeleteArticle(ctx context.Context, aid int) error {
	if s.deleteArticleFn != nil {
		return s.deleteArticleFn(ctx, aid)
	}
	return nil
}

func (s *indexerStub) EnsureBlogIndex() error {
	if s.ensureBlogIndexFn != nil {
		return s.ensureBlogIndexFn()
	}
	return nil
}

var _ meili.BlogSearchIndexer = (*indexerStub)(nil)

// --- Shared Stubs ---

type articleRepoStub struct {
	cacheFn             func() cache.Store
	getArticleFn        func(context.Context, int) (*blog.Article, error)
	listAllArticleIDsFn func(context.Context) ([]int, error)
	getLastAidFn        func(context.Context) (int, error)
	getContextFn        func(context.Context, int, int, int) (*blog.ArticleContext, error)
	getArticlesFn       func(context.Context, int, int) ([]blog.Article, error)
	getCountFn          func(context.Context) (int, error)

	getTagsFn            func(context.Context) ([]blog.Tag, error)
	getArticlesByTagFn   func(context.Context, string, int, int) ([]blog.Article, error)
	getArticlesByTagIDFn func(context.Context, int) ([]blog.ArticleDescription, error)
}

func (s *articleRepoStub) Cache() cache.Store {
	if s.cacheFn != nil {
		return s.cacheFn()
	}
	return &mockCache{}
}

func (s *articleRepoStub) GetArticle(ctx context.Context, aid int) (*blog.Article, error) {
	if s.getArticleFn != nil {
		return s.getArticleFn(ctx, aid)
	}
	return nil, nil
}

func (s *articleRepoStub) ListAllArticleIDs(ctx context.Context) ([]int, error) {
	if s.listAllArticleIDsFn != nil {
		return s.listAllArticleIDsFn(ctx)
	}
	return nil, nil
}

func (s *articleRepoStub) GetLastAid(ctx context.Context) (int, error) {
	if s.getLastAidFn != nil {
		return s.getLastAidFn(ctx)
	}
	return 0, nil
}

func (s *articleRepoStub) GetContext(ctx context.Context, startAid, direction, boundary int) (*blog.ArticleContext, error) {
	if s.getContextFn != nil {
		return s.getContextFn(ctx, startAid, direction, boundary)
	}
	return nil, nil
}

func (s *articleRepoStub) GetArticles(ctx context.Context, offset, limit int) ([]blog.Article, error) {
	if s.getArticlesFn != nil {
		return s.getArticlesFn(ctx, offset, limit)
	}
	return nil, nil
}

func (s *articleRepoStub) GetCount(ctx context.Context) (int, error) {
	if s.getCountFn != nil {
		return s.getCountFn(ctx)
	}
	return 0, nil
}

func (s *articleRepoStub) GetTags(ctx context.Context) ([]blog.Tag, error) {
	if s.getTagsFn != nil {
		return s.getTagsFn(ctx)
	}
	return nil, nil
}

func (s *articleRepoStub) GetArticlesByTag(ctx context.Context, tag string, offset, limit int) ([]blog.Article, error) {
	if s.getArticlesByTagFn != nil {
		return s.getArticlesByTagFn(ctx, tag, offset, limit)
	}
	return nil, nil
}

func (s *articleRepoStub) GetArticlesByTagID(ctx context.Context, id int) ([]blog.ArticleDescription, error) {
	if s.getArticlesByTagIDFn != nil {
		return s.getArticlesByTagIDFn(ctx, id)
	}
	return nil, nil
}

type accessRepoStub struct {
	getAccessFn     func(context.Context, int) ([]blog.Access, error)
	logAccessInfoFn func(context.Context, int, string) error
}

func (s *accessRepoStub) GetAccess(ctx context.Context, day int) ([]blog.Access, error) {
	if s.getAccessFn != nil {
		return s.getAccessFn(ctx, day)
	}
	return nil, nil
}

func (s *accessRepoStub) LogAccessInfo(ctx context.Context, aid int, ip string) error {
	if s.logAccessInfoFn != nil {
		return s.logAccessInfoFn(ctx, aid, ip)
	}
	return nil
}

type commentRepoStub struct {
	cacheFn                func() cache.Store
	getCommentsFn          func(context.Context, int, int, int) ([]blog.Comment, error)
	getCommentCountFn      func(context.Context, int) (int, error)
	getChildCommentsFn     func(context.Context, int, int, int) ([]blog.Comment, error)
	getChildCommentCountFn func(context.Context, int) (int, error)
	addCommentFn           func(context.Context, *blog.Comment) error
}

func (s *commentRepoStub) Cache() cache.Store {
	if s.cacheFn != nil {
		return s.cacheFn()
	}
	return &mockCache{}
}

func (s *commentRepoStub) GetComments(ctx context.Context, aid, offset, limit int) ([]blog.Comment, error) {
	if s.getCommentsFn != nil {
		return s.getCommentsFn(ctx, aid, offset, limit)
	}
	return nil, nil
}

func (s *commentRepoStub) GetCommentCount(ctx context.Context, aid int) (int, error) {
	if s.getCommentCountFn != nil {
		return s.getCommentCountFn(ctx, aid)
	}
	return 0, nil
}

func (s *commentRepoStub) GetChildComments(ctx context.Context, id, offset, limit int) ([]blog.Comment, error) {
	if s.getChildCommentsFn != nil {
		return s.getChildCommentsFn(ctx, id, offset, limit)
	}
	return nil, nil
}

func (s *commentRepoStub) GetChildCommentCount(ctx context.Context, id int) (int, error) {
	if s.getChildCommentCountFn != nil {
		return s.getChildCommentCountFn(ctx, id)
	}
	return 0, nil
}

func (s *commentRepoStub) AddComment(ctx context.Context, comment *blog.Comment) error {
	if s.addCommentFn != nil {
		return s.addCommentFn(ctx, comment)
	}
	return nil
}

// mockCache is a simple no-op cache for testing
type mockCache struct{}

func (m *mockCache) Get(ctx context.Context, key string) ([]byte, error) { return nil, nil }
func (m *mockCache) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return nil
}
func (m *mockCache) Del(ctx context.Context, key string) error          { return nil }
func (m *mockCache) DelMatch(ctx context.Context, pattern string) error { return nil }
func (m *mockCache) GetMatch(ctx context.Context, pattern string) (map[string][]byte, error) {
	return make(map[string][]byte), nil
}
func (m *mockCache) Flush(ctx context.Context) error { return nil }

var _ dao.ArticleRepository = (*articleRepoStub)(nil)
var _ dao.AccessRepository = (*accessRepoStub)(nil)
var _ dao.CommentRepository = (*commentRepoStub)(nil)
var _ dao.ManagementArticleRepository = (*managementRepoStub)(nil)
var _ dao.ManagementTagRepository = (*managementRepoStub)(nil)

type managementRepoStub struct {
	cacheFn                func() cache.Store
	createArticleFn        func(context.Context, *blog.Article) (int, error)
	CreateArticleWithAidFn func(context.Context, *blog.Article) error
	updateArticleFn        func(context.Context, int, *blog.ArticleUpdate) error
	deleteArticleFn        func(context.Context, int) error
	createTagFn            func(context.Context, string) error
	updateTagFn            func(context.Context, int, *blog.TagUpdate) error
	deleteTagFn            func(context.Context, int) error
}

func (s *managementRepoStub) Cache() cache.Store {
	if s.cacheFn != nil {
		return s.cacheFn()
	}
	return &mockCache{}
}

func (s *managementRepoStub) CreateArticle(ctx context.Context, article *blog.Article) (int, error) {
	if s.createArticleFn != nil {
		return s.createArticleFn(ctx, article)
	}
	return 0, nil
}

func (s *managementRepoStub) CreateArticleWithAid(ctx context.Context, article *blog.Article) error {
	if s.CreateArticleWithAidFn != nil {
		return s.CreateArticleWithAidFn(ctx, article)
	}
	return nil
}

func (s *managementRepoStub) UpdateArticle(ctx context.Context, aid int, updates *blog.ArticleUpdate) error {
	if s.updateArticleFn != nil {
		return s.updateArticleFn(ctx, aid, updates)
	}
	return nil
}

func (s *managementRepoStub) DeleteArticle(ctx context.Context, aid int) error {
	if s.deleteArticleFn != nil {
		return s.deleteArticleFn(ctx, aid)
	}
	return nil
}

func (s *managementRepoStub) CreateTag(ctx context.Context, name string) error {
	if s.createTagFn != nil {
		return s.createTagFn(ctx, name)
	}
	return nil
}

func (s *managementRepoStub) UpdateTag(ctx context.Context, id int, updates *blog.TagUpdate) error {
	if s.updateTagFn != nil {
		return s.updateTagFn(ctx, id, updates)
	}
	return nil
}

func (s *managementRepoStub) DeleteTag(ctx context.Context, id int) error {
	if s.deleteTagFn != nil {
		return s.deleteTagFn(ctx, id)
	}
	return nil
}
