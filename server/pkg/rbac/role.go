// Package rbac is the shared identity read model: role and permission records,
// cached reads, and cache-invalidation events. It is not a transport adapter.
//
// Request-path authorization stays in HTTP middleware (RequireAnyPermission).
// Cross-service identity on a request is authctx.TokenUser, not these records.
package rbac

import "github.com/ipuppet/aqua/server/pkg/types"

// Role is the shared role record used by the RBAC reader and user-service persistence.
type Role struct {
	ID          int              `db:"id" json:"id"`
	Code        string           `db:"code" json:"code"`
	Description types.NullString `db:"description" json:"description"`
}
