package logic

import (
	"context"

	"aqua/services/passport/dao"
	auth "aqua/services/passport/domain"
)

// SessionService is the default SessionService implementation.
type SessionService struct {
	repo dao.AuthRepository
}

// NewSessionService wires dependencies into a SessionService implementation.
func NewSessionService(repo dao.AuthRepository) *SessionService {
	return &SessionService{repo: repo}
}

// ListUserSessions retrieves active sessions for a given user from Redis-backed session index.
func (s *SessionService) ListUserSessions(ctx context.Context, uid int) ([]auth.Session, error) {
	return s.repo.GetUserSessions(ctx, uid)
}

// RevokeUserSession deletes a user session by its stable session id.
func (s *SessionService) RevokeUserSession(ctx context.Context, uid int, sessionID string) error {
	return s.repo.DeleteUserSession(ctx, uid, sessionID)
}
