package meili

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/meilisearch/meilisearch-go"

	"github.com/ipuppet/aqua/server/pkg/infra/config"
	"github.com/ipuppet/aqua/server/pkg/infra/logutil"
	"github.com/ipuppet/aqua/server/pkg/infra/models"
	blog "github.com/ipuppet/aqua/server/services/blog/domain"
)

const blogIndexName = "blog"

// BlogSearchIndexer defines the contract for blog article search and indexing.
type BlogSearchIndexer interface {
	SearchArticles(ctx context.Context, query string, limit, offset int) (*blog.SearchResponse, error)
	IndexArticle(ctx context.Context, article *blog.Article) error
	DeleteArticle(ctx context.Context, aid int) error
	EnsureBlogIndex() error
}

type blogSearchIndexer struct {
	meiliSearchConfigProvider config.MeiliSearchConfigProvider
}

// NewBlogSearchIndexer creates a BlogSearchIndexer from a runtime MeiliSearch config provider.
func NewBlogSearchIndexer(meiliSearchConfigProvider config.MeiliSearchConfigProvider) BlogSearchIndexer {
	if meiliSearchConfigProvider == nil {
		meiliSearchConfigProvider = func() models.MeiliSearchConfig {
			return models.MeiliSearchConfig{}
		}
	}
	return &blogSearchIndexer{meiliSearchConfigProvider: meiliSearchConfigProvider}
}

func (x *blogSearchIndexer) client() meilisearch.ServiceManager {
	if x == nil || x.meiliSearchConfigProvider == nil {
		return meilisearch.New("")
	}
	cfg := x.meiliSearchConfigProvider()
	return meilisearch.New(cfg.Host, meilisearch.WithAPIKey(cfg.APIKey))
}

func (x *blogSearchIndexer) EnsureBlogIndex() error {
	client := x.client()
	index := client.Index(blogIndexName)
	_, err := index.FetchInfo()
	if err != nil {
		resp, err := client.CreateIndex(&meilisearch.IndexConfig{
			Uid:        blogIndexName,
			PrimaryKey: "aid",
		})
		if err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
		slog.Info("Created MeiliSearch index", "task_uid", resp.TaskUID)
		_, _ = client.WaitForTask(resp.TaskUID, 0)
	}

	_, err = index.UpdateSortableAttributes(&[]string{"created_at_ts"})
	if err != nil {
		return fmt.Errorf("failed to update sortable attributes: %w", err)
	}

	filterable := []any{"tags"}
	_, err = index.UpdateFilterableAttributes(&filterable)
	if err != nil {
		return fmt.Errorf("failed to update filterable attributes: %w", err)
	}

	slog.Info("MeiliSearch blog index settings updated")
	return nil
}

func (x *blogSearchIndexer) IndexArticle(ctx context.Context, article *blog.Article) error {
	tags := make([]string, 0)
	if article.Tags != nil {
		for _, t := range article.Tags {
			tags = append(tags, t.Name)
		}
	}

	doc := ArticleDocument{
		Aid:     article.Aid,
		Title:   article.Title,
		Summary: article.Summary,
		Content: stripMarkdown(article.Markdown),
		Tags:    tags,
	}

	if article.CreatedAt != nil {
		doc.CreatedAt = article.CreatedAt.Format(time.RFC3339)
		doc.CreatedAtTs = article.CreatedAt.Unix()
	}

	task, err := x.client().Index(blogIndexName).AddDocuments([]ArticleDocument{doc}, nil)
	if err != nil {
		logutil.FromContext(ctx).Error("failed to send article to meilisearch",
			slog.Int("aid", article.Aid),
			slog.Any("error", err))
		return err
	}

	logutil.FromContext(ctx).Info("dispatched article indexing task",
		slog.Int("aid", article.Aid),
		slog.Int64("task_uid", task.TaskUID))
	return nil
}

func (x *blogSearchIndexer) DeleteArticle(ctx context.Context, aid int) error {
	task, err := x.client().Index(blogIndexName).DeleteDocument(fmt.Sprintf("%d", aid), nil)
	if err != nil {
		logutil.FromContext(ctx).Error("failed to delete article from meilisearch",
			slog.Int("aid", aid),
			slog.Any("error", err))
		return err
	}

	logutil.FromContext(ctx).Info("dispatched article deletion task",
		slog.Int("aid", aid),
		slog.Int64("task_uid", task.TaskUID))
	return nil
}

func (x *blogSearchIndexer) SearchArticles(ctx context.Context, query string, limit, offset int) (*blog.SearchResponse, error) {
	searchRes, err := x.client().Index(blogIndexName).Search(query, &meilisearch.SearchRequest{
		Limit:                 int64(limit),
		Offset:                int64(offset),
		Sort:                  []string{"created_at_ts:desc"},
		AttributesToHighlight: []string{"title", "summary", "content"},
		HighlightPreTag:       "<em>",
		HighlightPostTag:      "</em>",
		AttributesToCrop:      []string{"content"},
		CropLength:            50,
		CropMarker:            "...",
	})
	if err != nil {
		return nil, err
	}

	articles := make([]blog.Article, 0, len(searchRes.Hits))
	for _, hit := range searchRes.Hits {
		var aidFloat float64
		var title, summary, createdAtStr string
		var tags []string

		if v, ok := hit["aid"]; ok {
			_ = json.Unmarshal(v, &aidFloat)
		}
		if v, ok := hit["title"]; ok {
			_ = json.Unmarshal(v, &title)
		}
		if v, ok := hit["summary"]; ok {
			_ = json.Unmarshal(v, &summary)
		}
		if v, ok := hit["created_at"]; ok {
			_ = json.Unmarshal(v, &createdAtStr)
		}
		if v, ok := hit["tags"]; ok {
			_ = json.Unmarshal(v, &tags)
		}

		var createdAt *time.Time
		if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
			createdAt = &t
		}

		articleTags := make([]blog.Tag, 0, len(tags))
		for _, t := range tags {
			articleTags = append(articleTags, blog.Tag{Name: t})
		}

		article := &blog.Article{
			Aid:       int(aidFloat),
			Title:     title,
			Summary:   summary,
			Tags:      articleTags,
			CreatedAt: createdAt,
		}

		if formattedBytes, ok := hit["_formatted"]; ok {
			var formattedMap map[string]any
			if err := json.Unmarshal(formattedBytes, &formattedMap); err == nil {
				if v, ok := formattedMap["title"]; ok {
					article.Title = fmt.Sprintf("%v", v)
				}

				var hlSummary, hlContent string
				if v, ok := formattedMap["summary"]; ok {
					hlSummary = fmt.Sprintf("%v", v)
				}
				if v, ok := formattedMap["content"]; ok {
					hlContent = fmt.Sprintf("%v", v)
				}

				if hlSummary != "" && strings.Contains(hlSummary, "<em>") {
					article.Summary = hlSummary
				} else if hlContent != "" && strings.Contains(hlContent, "<em>") {
					article.Summary = hlContent
				}
			}
		}

		articles = append(articles, *article)
	}

	return &blog.SearchResponse{
		Items: articles,
		Total: int(searchRes.EstimatedTotalHits),
	}, nil
}
