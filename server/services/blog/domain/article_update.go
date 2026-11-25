package blog

// ArticleUpdate represents the modifiable properties of an article.
// Pointer types are used to distinguish between intentional omissions (nil) and actual zero values during partial updates.
type ArticleUpdate struct {
	Title    *string `json:"title" db:"title"`
	Summary  *string `json:"summary" db:"summary"`
	Markdown *string `json:"markdown" db:"markdown"`
	// Tags requires custom processing logic to handle array updates.
	Tags *[]Tag `json:"tags"`
}
