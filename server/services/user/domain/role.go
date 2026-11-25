package user

// RoleResponse is the user-service API view of a role.
type RoleResponse struct {
	ID          int    `json:"id"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

// RolesResponse represents a list of roles.
type RolesResponse struct {
	Roles []RoleResponse `json:"roles"`
}
