package grpcserver

//go:generate go tool goverter gen .

import (
	passportv1 "github.com/ipuppet/aqua/server/gen/go/aqua/passport/v1"
	auth "github.com/ipuppet/aqua/server/services/passport/domain"
)

// goverter:variables
// goverter:output:file ./generated.go
// goverter:extend github.com/ipuppet/aqua/server/pkg/utils/convert:ToPBTime
// goverter:extend github.com/ipuppet/aqua/server/pkg/utils/convert:IntToInt32
// goverter:useZeroValueOnPointerInconsistency
// goverter:ignoreUnexported
var (
	// goverter:map UID Uid
	// goverter:map SessionID SessionId
	ToUserSession  func(item auth.Session) *passportv1.UserSession
	ToUserSessions func(items []auth.Session) []*passportv1.UserSession
)
