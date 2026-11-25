package dao

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	auth "github.com/ipuppet/aqua/server/services/passport/domain"
)

func userSessionKeyPattern(uid int) string {
	return "session:" + strconv.Itoa(uid) + ":*"
}

func sessionKey(uid int, sessionID string) string {
	return "session:" + strconv.Itoa(uid) + ":" + sessionID
}

// StoreUserSession persists session metadata by stable session id with the given TTL.
func (r *repository) StoreUserSession(ctx context.Context, session *auth.Session, ttl time.Duration) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	b, err := json.Marshal(session)
	if err != nil {
		return err
	}
	return r.cache.Set(ctx, sessionKey(session.UID, session.SessionID), b, ttl)
}

// GetUserSession returns session metadata bound to the given session id if it exists.
func (r *repository) GetUserSession(ctx context.Context, uid int, sessionID string) (*auth.Session, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	val, err := r.cache.Get(ctx, sessionKey(uid, sessionID))
	if err != nil {
		return nil, err
	}
	var session auth.Session
	err = json.Unmarshal(val, &session)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

// GetUserSessions returns recent sessions for a user by querying all sessions matching the user pattern.
func (r *repository) GetUserSessions(ctx context.Context, uid int) ([]auth.Session, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if uid <= 0 {
		return []auth.Session{}, nil
	}

	// Get all sessions for this user using pattern matching
	sessionMap, err := r.cache.GetMatch(ctx, userSessionKeyPattern(uid))
	if err != nil {
		return nil, err
	}

	if len(sessionMap) == 0 {
		return []auth.Session{}, nil
	}

	out := make([]auth.Session, 0, len(sessionMap))
	for _, data := range sessionMap {
		var session auth.Session
		if err := json.Unmarshal(data, &session); err != nil {
			continue
		}
		out = append(out, session)
	}

	return out, nil
}

// DeleteUserSession removes a session from Redis (e.g. on rotation or logout).
func (r *repository) DeleteUserSession(ctx context.Context, uid int, sessionID string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return r.cache.Del(ctx, sessionKey(uid, sessionID))
}

// DeleteAllUserSessions removes all user sessions from Redis.
func (r *repository) DeleteAllUserSessions(ctx context.Context, uid int) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if uid <= 0 {
		return nil
	}
	return r.cache.DelMatch(ctx, userSessionKeyPattern(uid))
}
