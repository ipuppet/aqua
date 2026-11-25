package grpcserver

import (
	"context"

	passportv1 "aqua/gen/go/aqua/passport/v1"
	"aqua/services/passport/logic"
)

// AdminServer provides gRPC APIs for admin-side passport operations.
type AdminServer struct {
	sessionService *logic.SessionService
	passportv1.UnimplementedAdminPassportServiceServer
}

// NewAdminServer constructs AdminServer with the required dependencies.
func NewAdminServer(sessionService *logic.SessionService) *AdminServer {
	return &AdminServer{sessionService: sessionService}
}

// ListUserSessions returns all active sessions for a given user.
func (s *AdminServer) ListUserSessions(ctx context.Context, req *passportv1.ListUserSessionsRequest) (*passportv1.ListUserSessionsResponse, error) {
	sessions, err := s.sessionService.ListUserSessions(ctx, int(req.GetUid()))
	if err != nil {
		return nil, err
	}
	return &passportv1.ListUserSessionsResponse{Sessions: ToUserSessions(sessions)}, nil
}

// RevokeUserSession revokes a session by session_id.
func (s *AdminServer) RevokeUserSession(ctx context.Context, req *passportv1.RevokeUserSessionRequest) (*passportv1.RevokeUserSessionResponse, error) {
	sessionID := req.GetSessionId()
	if err := s.sessionService.RevokeUserSession(ctx, int(req.GetUid()), sessionID); err != nil {
		return nil, err
	}
	return &passportv1.RevokeUserSessionResponse{}, nil
}
