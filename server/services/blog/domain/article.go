package blog

import "time"

// Tag represents a categorization tag associated with an article.
type Tag struct {
	ID   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

// Article represents a complete blog article including its metadata and full content.
type Article struct {
	Aid       int        `db:"aid" json:"aid"`
	Title     string     `db:"title" json:"title"`
	Tags      []Tag      `json:"tags"`
	Summary   string     `db:"summary" json:"summary"`
	Markdown  string     `db:"markdown" json:"markdown"`
	CreatedAt *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
}

// ArticleContext represents a minimal reference to an article, typically used in navigational links.
type ArticleContext struct {
	Aid   int    `db:"aid" json:"aid"`
	Title string `db:"title" json:"title"`
}

// ArticleDescription represents a summary view of an article, excluding its full body content.
type ArticleDescription struct {
	Aid       int        `db:"aid" json:"aid"`
	Title     string     `db:"title" json:"title"`
	Tags      []Tag      `json:"tags"`
	Summary   string     `db:"summary" json:"summary"`
	CreatedAt *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
}

// ArticlesResponse represents a paginated list of articles alongside the dataset's total count.
type ArticlesResponse struct {
	Articles []Article `json:"articles"`
	Total    int       `json:"total"`
}

// SearchResponse represents the standard payload for search results and overall match count.
type SearchResponse struct {
	Items []Article `json:"items"`
	Total int       `json:"total"`
}

// ArticleDetail is a public article view plus previous/next context.
type ArticleDetail struct {
	Article  *Article
	Previous *ArticleContext
	Next     *ArticleContext
}
