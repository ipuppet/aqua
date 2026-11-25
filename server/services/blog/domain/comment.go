package blog

import "time"

// Comment represents a user-submitted comment on an article.
type Comment struct {
	ID        int        `db:"id" json:"id"`
	UID       int        `db:"uid" json:"uid"`
	Aid       int        `db:"aid" json:"aid"`
	Content   string     `db:"content" json:"content"`
	Parent    int        `db:"parent" json:"parent"`
	IP        string     `db:"ip" json:"ip"`
	CreatedAt *time.Time `db:"created_at" json:"created_at"`
	IsDeleted bool       `db:"is_deleted" json:"is_deleted"`
	ReplyTo   string     `db:"reply_to" json:"reply_to"`
}

// CommentPage is a paginated comment list.
type CommentPage struct {
	Comments []Comment
	Total    int
}
