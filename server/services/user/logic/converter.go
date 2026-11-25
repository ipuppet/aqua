package logic

//go:generate go tool goverter gen .

import (
	"aqua/pkg/rbac"
	"aqua/pkg/types"
	user "aqua/services/user/domain"
)

func nullStringToString(value types.NullString) string {
	if !value.Valid {
		return ""
	}
	return value.String
}

// goverter:variables
// goverter:output:file ./generated.go
// goverter:extend aqua/services/user/logic:nullStringToString
// goverter:extend aqua/pkg/utils/convert:CloneTime
// goverter:useZeroValueOnPointerInconsistency
// goverter:ignoreUnexported
var (
	ToUserModelResponse  func(item *user.UserModel) *user.UserModelResponse
	ToUserModelResponses func(items []user.UserModel) []user.UserModelResponse

	ToRoleResponse       func(item *rbac.Role) *user.RoleResponse
	ToPermissionResponse func(item *rbac.Permission) *user.PermissionResponse
)
