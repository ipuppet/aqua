package grpcserver

//go:generate go tool goverter gen .

import (
	passportv1 "aqua/gen/go/aqua/passport/v1"
	auth "aqua/services/passport/domain"
)

// goverter:variables
// goverter:output:file ./generated.go
// goverter:extend aqua/pkg/utils/convert:ToPBTime
// goverter:extend aqua/pkg/utils/convert:IntToInt32
// goverter:useZeroValueOnPointerInconsistency
// goverter:ignoreUnexported
var (
	// goverter:map UID Uid
	// goverter:map SessionID SessionId
	ToUserSession  func(item auth.Session) *passportv1.UserSession
	ToUserSessions func(items []auth.Session) []*passportv1.UserSession
)
