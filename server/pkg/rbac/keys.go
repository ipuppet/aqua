package rbac

import "strconv"

const (
	userInfoKeyPre  = "user:GetInfo:"
	userRolesKeyPre = "user:GetRoles:"
	rolePermsKeyPre = "role:GetPermissions:"
)

const (
	UserRolesPattern       = "user:GetRoles:*"
	RolePermissionsPattern = "role:GetPermissions:*"
)

// CacheKeyUserInfo returns the cache key for user info.
func CacheKeyUserInfo(uid int) string {
	return userInfoKeyPre + strconv.Itoa(uid)
}

// CacheKeyUserRoles returns the cache key for user-role mapping.
func CacheKeyUserRoles(uid int) string {
	return userRolesKeyPre + strconv.Itoa(uid)
}

// CacheKeyStringRolePermissions returns the cache key for a role's permissions.
func CacheKeyStringRolePermissions(role string) string {
	return rolePermsKeyPre + role
}
