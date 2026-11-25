package user

// PermissionUpdate represents the modifiable properties of a permission.
// Pointer types are used to distinguish between intentional omissions (nil) and actual zero values during updates.
type PermissionUpdate struct {
	ID          *int    `json:"id" db:"-"`
	Code        *string `json:"code" db:"code"`
	Description *string `json:"description" db:"description"`
}
