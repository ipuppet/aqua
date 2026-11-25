package logic

//go:generate go tool goverter gen .

import (
	"github.com/ipuppet/aqua/server/pkg/rbac"
	"github.com/ipuppet/aqua/server/pkg/types"
	user "github.com/ipuppet/aqua/server/services/user/domain"
)

func nullStringToString(value types.NullString) string {
	if !value.Valid {
		return ""
	}
	return value.String
}

// goverter:variables
// goverter:output:file ./generated.go
// goverter:extend github.com/ipuppet/aqua/server/services/user/logic:nullStringToString
// goverter:extend github.com/ipuppet/aqua/server/pkg/utils/convert:CloneTime
// goverter:useZeroValueOnPointerInconsistency
// goverter:ignoreUnexported
var (
	ToUserModelResponse  func(item *user.UserModel) *user.UserModelResponse
	ToUserModelResponses func(items []user.UserModel) []user.UserModelResponse

	ToRoleResponse       func(item *rbac.Role) *user.RoleResponse
	ToPermissionResponse func(item *rbac.Permission) *user.PermissionResponse
)
