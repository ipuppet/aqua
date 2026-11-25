package rbac

import "github.com/ipuppet/aqua/server/pkg/types"

// Permission is the shared permission record used by the RBAC reader and user-service persistence.
type Permission struct {
	ID          int              `db:"id" json:"id"`
	Code        string           `db:"code" json:"code"`
	Description types.NullString `db:"description" json:"description"`
}
