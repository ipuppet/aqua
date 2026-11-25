package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"aqua/pkg/infra/database"
	blog "aqua/services/blog/domain"
)

func (r *repository) GetArticle(ctx context.Context, aid int) (*blog.Article, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	article := &blog.Article{}
	err := database.FetchWithCache(
		ctx,
		r.cache,
		"blog:GetArticle:"+strconv.Itoa(aid),
		24*time.Hour,
		article,
		func() (any, error) {
			a := &blog.Article{}
			err := database.Retrieve(
				ctx,
				r.db,
				a,
				`select
				aid,
				title,
				summary,
				markdown,
				created_at,
				updated_at
				from blog.article
				where aid=$1
				limit 1`,
				aid)
			return a, err
		})
	if err != nil {
		return nil, fmt.Errorf("getting article with aid %d: %w", aid, err)
	}

	tags := []blog.Tag{}
	err = database.FetchWithCache(
		ctx,
		r.cache,
		"blog:GetArticleTags:"+strconv.Itoa(aid),
		24*time.Hour,
		&tags,
		func() (any, error) {
			t := []blog.Tag{}
			err := database.Retrieve(
				ctx,
				r.db,
				&t,
				`select
				t.id, t.name
				from blog.tag t
				inner join blog.tag_article ta on t.id = ta.tag_id
				where ta.aid = $1
				order by t.id`,
				aid)
			return t, err
		})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("getting tags for article with aid %d: %w", aid, err)
	}

	article.Tags = tags
	return article, nil
}

func (r *repository) ListAllArticleIDs(ctx context.Context) ([]int, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	var rows []struct {
		Aid int `db:"aid"`
	}
	err := database.Retrieve(ctx, r.db, &rows, "SELECT aid FROM blog.article ORDER BY aid")
	if err != nil {
		return nil, fmt.Errorf("listing article ids: %w", err)
	}
	ids := make([]int, len(rows))
	for i, row := range rows {
		ids[i] = row.Aid
	}
	return ids, nil
}

func (r *repository) GetCount(ctx context.Context) (int, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}

	s := &struct {
		Count int `db:"count"`
	}{}
	err := database.FetchWithCache(
		ctx,
		r.cache,
		"blog:GetCount",
		24*time.Hour,
		s,
		func() (any, error) {
			res := &struct {
				Count int `db:"count"`
			}{}
			err := database.Retrieve(ctx, r.db, res, "select count(aid) as count from blog.article")
			return res, err
		})
	if err != nil {
		return 0, fmt.Errorf("getting article count: %w", err)
	}
	return s.Count, nil
}

func (r *repository) GetLastAid(ctx context.Context) (int, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}

	s := &struct {
		Aid int `db:"aid"`
	}{}
	err := database.FetchWithCache(
		ctx,
		r.cache,
		"blog:GetLastAid",
		24*time.Hour,
		s,
		func() (any, error) {
			res := &struct {
				Aid int `db:"aid"`
			}{}
			err := database.Retrieve(ctx, r.db, res, "select aid from blog.article order by aid desc limit 1")
			return res, err
		})
	if err != nil {
		return 0, fmt.Errorf("getting last aid: %w", err)
	}
	return s.Aid, nil
}

// GetContext finds the next or previous article.
// direction: 1 for next (aid increases), -1 for previous (aid decreases).
// boundaryAid: for direction 1, this is the max aid; for direction -1, this is effectively 0 or min aid.
// Returns (nil, nil) if no context article is found within bounds.
func (r *repository) GetContext(ctx context.Context, startAid int, direction int, lastKnownAid int) (*blog.ArticleContext, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	currentAid := startAid

	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		targetAid := currentAid + direction

		if direction > 0 && targetAid > lastKnownAid { // Moving forward, exceeded last known aid
			return nil, nil
		}
		if direction < 0 && targetAid <= 0 { // Moving backward, reached or passed first possible aid
			return nil, nil
		}

		ac := &blog.ArticleContext{}
		cacheKey := "blog:GetContext:" + strconv.Itoa(targetAid)
		err := database.FetchWithCache(
			ctx,
			r.cache,
			cacheKey,
			24*time.Hour,
			ac,
			func() (any, error) {
				res := &blog.ArticleContext{}
				err := database.Retrieve(
					ctx,
					r.db,
					res,
					"SELECT aid, title FROM blog.article WHERE aid=$1 LIMIT 1", targetAid)
				return res, err
			})
		if err == nil {
			if ac.Aid == targetAid { // Found and populated
				return ac, nil
			}
			// If err is nil but ac not populated as expected (e.g. ac.Aid is 0),
			// it implies "not found" by the query. Continue searching.
			currentAid = targetAid
			continue
		}

		if errors.Is(err, sql.ErrNoRows) { // Explicitly not found
			currentAid = targetAid
			continue // Try next/previous
		}

		// Any other error
		return nil, fmt.Errorf("getting context for target aid %d (started from %d, direction %d): %w", targetAid, startAid, direction, err)
	}
}

func (r *repository) GetArticles(ctx context.Context, offset int, limit int) ([]blog.Article, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	articles := []blog.Article{}
	err := database.FetchWithCache(
		ctx,
		r.cache,
		"blog:GetArticles:"+strconv.Itoa(offset)+":"+strconv.Itoa(limit),
		24*time.Hour,
		&articles,
		func() (any, error) {
			res := []blog.Article{}
			err := database.Retrieve(
				ctx,
				r.db,
				&res,
				`select
				aid,
				title,
				summary,
				created_at,
				updated_at
				from blog.article
				order by created_at desc
				limit $1 offset $2`,
				limit, offset)
			return res, err
		})
	if err != nil {
		return nil, fmt.Errorf("getting articles with offset %d, limit %d: %w", offset, limit, err)
	}

	err = attachTagsToArticles(ctx, r.db, articles)
	if err != nil {
		return nil, err
	}

	return articles, nil
}

// GetArticlesByTagID gets articles by tag ID.
func (r *repository) GetArticlesByTagID(ctx context.Context, id int) ([]blog.ArticleDescription, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	articles := []blog.ArticleDescription{}
	err := database.Retrieve(
		ctx,
		r.db,
		&articles,
		`select distinct
			article.aid,
			article.title,
			article.summary,
			article.created_at,
			article.updated_at
		from blog.article
		inner join blog.tag_article on article.aid = tag_article.aid
		where tag_article.tag_id = $1
		order by article.created_at desc`,
		id)
	if err != nil {
		return nil, fmt.Errorf("getting articles for tag_id %d: %w", id, err)
	}

	// Attach tags to articles
	if err := attachTagsToArticleDescriptions(ctx, r.db, articles); err != nil {
		return nil, fmt.Errorf("attaching tags to articles for tag_id %d: %w", id, err)
	}

	return articles, nil
}

func (r *repository) GetTags(ctx context.Context) ([]blog.Tag, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	tags := []blog.Tag{}
	err := database.FetchWithCache(
		ctx,
		r.cache,
		"blog:GetTags",
		24*time.Hour,
		&tags,
		func() (any, error) {
			res := []blog.Tag{}
			err := database.Retrieve(
				ctx,
				r.db,
				&res,
				"select * from blog.tag order by id")
			return res, err
		})
	if err != nil {
		return nil, fmt.Errorf("getting tags: %w", err)
	}
	return tags, err
}

// Helper function to get tags for multiple articles
func getTagsForArticles(ctx context.Context, db database.QueryContextExecutor, articleIDs []int) (map[int][]blog.Tag, error) {
	if len(articleIDs) == 0 {
		return make(map[int][]blog.Tag), nil
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	// Build placeholders for IN clause
	placeholders := make([]string, len(articleIDs))
	args := make([]any, len(articleIDs))
	for i, id := range articleIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	var allTags []*struct {
		Aid  int    `db:"aid"`
		ID   int    `db:"id"`
		Name string `db:"name"`
	}

	query := fmt.Sprintf(`select
		ta.aid,
		t.id,
		t.name
		from blog.tag t
		inner join blog.tag_article ta on t.id = ta.tag_id
		where ta.aid in (%s)
		order by ta.aid, t.id`,
		strings.Join(placeholders, ","))

	err := database.Retrieve(ctx, db, &allTags, query, args...)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("getting tags for articles: %w", err)
	}

	// Group tags by article ID
	tagsByAid := make(map[int][]blog.Tag)

	// Initialize empty arrays for all articles (handles articles with no tags)
	for _, aid := range articleIDs {
		tagsByAid[aid] = []blog.Tag{}
	}

	// Fill in the tags for articles that have them
	for _, tagRow := range allTags {
		tagsByAid[tagRow.Aid] = append(tagsByAid[tagRow.Aid], blog.Tag{
			ID:   tagRow.ID,
			Name: tagRow.Name,
		})
	}

	return tagsByAid, nil
}

// Helper function to attach tags to articles - reduces code duplication
func attachTagsToArticles(ctx context.Context, db database.QueryContextExecutor, articles []blog.Article) error {
	if len(articles) == 0 {
		return nil
	}

	if err := ctx.Err(); err != nil {
		return err
	}

	articleIDs := make([]int, len(articles))
	for i, article := range articles {
		articleIDs[i] = article.Aid
	}

	tagsByAid, err := getTagsForArticles(ctx, db, articleIDs)
	if err != nil {
		return err
	}

	// Assign tags to articles
	for i, article := range articles {
		articles[i].Tags = tagsByAid[article.Aid]
	}

	return nil
}

// Helper function to attach tags to article descriptions
func attachTagsToArticleDescriptions(ctx context.Context, db database.QueryContextExecutor, articles []blog.ArticleDescription) error {
	if len(articles) == 0 {
		return nil
	}

	if err := ctx.Err(); err != nil {
		return err
	}

	articleIDs := make([]int, len(articles))
	for i, article := range articles {
		articleIDs[i] = article.Aid
	}

	tagsByAid, err := getTagsForArticles(ctx, db, articleIDs)
	if err != nil {
		return err
	}

	// Assign tags to article descriptions
	for i, article := range articles {
		articles[i].Tags = tagsByAid[article.Aid]
	}

	return nil
}
