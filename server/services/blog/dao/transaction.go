package dao

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	blog "github.com/ipuppet/aqua/server/services/blog/domain"
)

// executeInTransaction runs callback within a transaction respecting context cancellation.
func (r *repository) executeInTransaction(ctx context.Context, fn func(context.Context, *sql.Tx) error) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}

	if err := fn(ctx, tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("%w; rollback: %v", err, rbErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}

// connectAidAndTags associates an article with tags within a transaction.
func connectAidAndTags(aid int, tags []blog.Tag) (string, []any) {
	if len(tags) == 0 {
		return "", nil
	}

	var valueStrings []string
	var args []any
	for i, tag := range tags {
		base := i * 2
		valueStrings = append(valueStrings, fmt.Sprintf("($%d,$%d)", base+1, base+2))
		args = append(args, aid, tag.ID)
	}

	return strings.Join(valueStrings, ","), args
}

// insertTagAssociations inserts tag associations for an article within a transaction.
func insertTagAssociations(ctx context.Context, tx *sql.Tx, aid int, tags []blog.Tag) error {
	if len(tags) == 0 {
		return nil
	}

	tagsValuesSQL, args := connectAidAndTags(aid, tags)
	if tagsValuesSQL == "" {
		return nil
	}

	_, err := tx.ExecContext(ctx, "insert into blog.tag_article (aid,tag_id) values "+tagsValuesSQL, args...)
	if err != nil {
		return fmt.Errorf("inserting tag associations for article %d: %w", aid, err)
	}
	return nil
}
