package dao

import (
	"context"
	"strconv"
	"time"

	"github.com/ipuppet/aqua/server/pkg/infra/database"
	blog "github.com/ipuppet/aqua/server/services/blog/domain"
)

func (r *repository) GetComments(ctx context.Context, aid int, offset int, limit int) ([]blog.Comment, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	result := []blog.Comment{}
	err := database.Retrieve(
		ctx,
		r.db,
		&result,
		`select c.id, c.uid, c.aid, c.content, c.parent, c.ip, c.created_at, c.is_deleted, c.reply_to
		from blog.comment c
		where c.is_deleted=false
		and c.aid=$1 and c.parent=-1
		order by c.created_at desc
		limit $2 offset $3`,
		aid, limit, offset)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *repository) GetCommentCount(ctx context.Context, aid int) (int, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}

	s := &struct {
		Count int `db:"count"`
	}{}
	err := database.FetchWithCache(
		ctx,
		r.cache,
		"blog:GetCommentCount:"+strconv.Itoa(aid),
		24*time.Hour,
		s,
		func() (any, error) {
			res := &struct {
				Count int `db:"count"`
			}{}
			err := database.Retrieve(ctx, r.db, res, "select count(id) as count from blog.comment where aid=$1 and parent=-1", aid)
			return res, err
		})
	if err != nil {
		return -1, err
	}

	return s.Count, nil
}

func (r *repository) GetChildComments(ctx context.Context, id int, offset int, limit int) ([]blog.Comment, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	result := []blog.Comment{}
	err := database.Retrieve(
		ctx,
		r.db,
		&result,
		`select c.id, c.uid, c.aid, c.content, c.parent, c.ip, c.created_at, c.is_deleted, c.reply_to
		from blog.comment c
		where c.is_deleted=false
		and c.parent=$1
		order by c.created_at desc
		limit $2 offset $3`,
		id, limit, offset)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *repository) GetChildCommentCount(ctx context.Context, id int) (int, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}

	s := &struct {
		Count int `db:"count"`
	}{}
	err := database.FetchWithCache(
		ctx,
		r.cache,
		"blog:GetChildCommentCount:"+strconv.Itoa(id),
		24*time.Hour,
		s,
		func() (any, error) {
			res := &struct {
				Count int `db:"count"`
			}{}
			err := database.Retrieve(ctx, r.db, res, "select count(id) as count from blog.comment where parent=$1", id)
			return res, err
		})
	if err != nil {
		return -1, err
	}

	return s.Count, nil
}

func (r *repository) AddComment(ctx context.Context, comment *blog.Comment) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	_, err := r.db.ExecContext(ctx, `
		insert into blog.comment (uid, aid, content, parent, ip, reply_to)
		values ($1, $2, $3, $4, $5, $6)`,
		comment.UID,
		comment.Aid,
		comment.Content,
		comment.Parent,
		comment.IP,
		comment.ReplyTo,
	)
	if err != nil {
		return err
	}

	if r.cache != nil {
		_ = r.cache.Del(ctx, "blog:GetCommentCount:"+strconv.Itoa(comment.Aid))
		_ = r.cache.Del(ctx, "blog:GetChildCommentCount:"+strconv.Itoa(comment.Parent))
	}

	return nil
}
