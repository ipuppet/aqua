package blog

// TagUpdate represents the modifiable properties of a tag.
// Pointer types are used to distinguish between intentional omissions (nil) and actual zero values during partial updates.
type TagUpdate struct {
	Name *string `json:"name" db:"name"`
}
