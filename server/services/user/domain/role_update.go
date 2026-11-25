package user

// RoleUpdate represents the modifiable properties of a role.
// Pointer types are used to distinguish between intentional omissions (nil) and actual zero values during updates.
type RoleUpdate struct {
	ID          *int    `json:"id" db:"-"`
	Code        *string `json:"code" db:"code"`
	Description *string `json:"description" db:"description"`
}
