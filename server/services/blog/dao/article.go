package dao

import (
	"context"
	"database/sql"
	"fmt"

	"aqua/pkg/infra/database"
	blog "aqua/services/blog/domain"
)

// CreateArticle creates a new article with tags.
func (r *repository) CreateArticle(ctx context.Context, article *blog.Article) (int, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}

	var aid int
	err := r.executeInTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		err := tx.QueryRowContext(ctx, "insert into blog.article (title, summary, markdown) values ($1, $2, $3) RETURNING aid", article.Title, article.Summary, article.Markdown).Scan(&aid)
		if err != nil {
			return fmt.Errorf("inserting article: %w", err)
		}
		return insertTagAssociations(ctx, tx, aid, article.Tags)
	})
	if err != nil {
		return 0, err
	}
	return aid, nil
}

// CreateArticleWithAid creates an article with specified aid (for special articles like about page).
func (r *repository) CreateArticleWithAid(ctx context.Context, article *blog.Article) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	return r.executeInTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, "insert into blog.article (aid, title, summary, markdown, created_at) values ($1, $2, $3, $4, $5)", article.Aid, article.Title, article.Summary, article.Markdown, article.CreatedAt)
		if err != nil {
			return fmt.Errorf("inserting article with aid %d: %w", article.Aid, err)
		}
		return insertTagAssociations(ctx, tx, article.Aid, article.Tags)
	})
}

// UpdateArticle updates an article and optionally replaces associated tags.
func (r *repository) UpdateArticle(ctx context.Context, aid int, updates *blog.ArticleUpdate) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	return r.executeInTransaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		if updates.Title != nil || updates.Summary != nil || updates.Markdown != nil {
			err := database.ExecuteUpdateTx(tx, "blog.article", updates, "aid=$1", aid)
			if err != nil {
				return fmt.Errorf("updating article %d: %w", aid, err)
			}
		}

		if updates.Tags == nil {
			return nil
		}

		_, err := tx.ExecContext(ctx, "delete from blog.tag_article where aid=$1", aid)
		if err != nil {
			return fmt.Errorf("deleting old tag associations for article %d: %w", aid, err)
		}
		return insertTagAssociations(ctx, tx, aid, *updates.Tags)
	})
}

// DeleteArticle deletes an article.
func (r *repository) DeleteArticle(ctx context.Context, aid int) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	_, err := r.db.ExecContext(ctx, "delete from blog.article where aid=$1", aid)
	if err != nil {
		return fmt.Errorf("deleting article %d: %w", aid, err)
	}
	return nil
}
