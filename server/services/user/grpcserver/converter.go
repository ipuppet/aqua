package grpcserver

//go:generate go tool goverter gen .

import (
	userv1 "github.com/ipuppet/aqua/server/gen/go/aqua/user/v1"
	user "github.com/ipuppet/aqua/server/services/user/domain"
)

// goverter:variables
// goverter:output:file ./generated.go
// goverter:extend aqua/pkg/utils/convert:ToPBTime
// goverter:extend aqua/pkg/utils/convert:IntToInt32
// goverter:useZeroValueOnPointerInconsistency
// goverter:ignoreUnexported
var (
	// goverter:map UID Uid
	ToPBUser  func(item user.UserModelResponse) *userv1.User
	ToPBUsers func(items []user.UserModelResponse) []*userv1.User

	// goverter:map ID Id
	ToPBRole  func(item user.RoleResponse) *userv1.Role
	ToPBRoles func(items []user.RoleResponse) []*userv1.Role
	// goverter:map ID Id
	ToPBPermission  func(item user.PermissionResponse) *userv1.Permission
	ToPBPermissions func(items []user.PermissionResponse) []*userv1.Permission
)
