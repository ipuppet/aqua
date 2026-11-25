package user

// UserUpdate represents the modifiable properties of a user.
// Pointer types are used to distinguish between intentional omissions (nil) and actual zero values.
// Note: Username is immutable after registration.
type UserUpdate struct {
	Nickname *string `json:"nickname" db:"nickname"`
	Email    *string `json:"email" db:"email"`
	Phone    *string `json:"phone" db:"phone"`
	Avatar   *string `json:"avatar" db:"avatar"`
	Gender   *int    `json:"gender" db:"gender"`
	Lock     *int    `json:"lock" db:"lock"`
}
