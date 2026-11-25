package dao

import (
	"context"
	"fmt"

	"github.com/ipuppet/aqua/server/pkg/infra/database"
	blog "github.com/ipuppet/aqua/server/services/blog/domain"
)

// CreateTag creates a new tag.
func (r *repository) CreateTag(ctx context.Context, name string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	_, err := r.db.ExecContext(ctx, "insert into blog.tag (name) values ($1)", name)
	if err != nil {
		return fmt.Errorf("adding tag '%s': %w", name, err)
	}
	return nil
}

// UpdateTag updates a tag.
func (r *repository) UpdateTag(ctx context.Context, id int, updates *blog.TagUpdate) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	query, args, err := database.BuildUpdateQuery("blog.tag", updates, "id=$1", id)
	if err != nil {
		return fmt.Errorf("building update for tag %d: %w", id, err)
	}

	if _, err := r.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("updating tag %d: %w", id, err)
	}
	return nil
}

// DeleteTag deletes a tag.
func (r *repository) DeleteTag(ctx context.Context, id int) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	_, err := r.db.ExecContext(ctx, "delete from blog.tag where id=$1", id)
	if err != nil {
		return fmt.Errorf("deleting tag %d: %w", id, err)
	}
	return nil
}
