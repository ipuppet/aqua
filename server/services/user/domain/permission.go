package user

// PermissionResponse is the user-service API view of a permission.
type PermissionResponse struct {
	ID          int    `json:"id"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

// PermissionsResponse represents a list of permissions.
type PermissionsResponse struct {
	Permissions []PermissionResponse `json:"permissions"`
}
