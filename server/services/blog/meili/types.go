package meili

// ArticleDocument represents the data structure stored in MeiliSearch
type ArticleDocument struct {
	Aid         int      `json:"aid"`
	Title       string   `json:"title"`
	Summary     string   `json:"summary"`
	Content     string   `json:"content"` // Plain text content stripped from Markdown
	Tags        []string `json:"tags"`
	CreatedAt   string   `json:"created_at"`
	CreatedAtTs int64    `json:"created_at_ts"` // For sorting
}
