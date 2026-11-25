package logic

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"aqua/pkg/infra/cache"
	"aqua/pkg/rbac"
	"aqua/services/passport/dao"
	auth "aqua/services/passport/domain"
)

func mustTestPrivateKey(t *testing.T) *rsa.PrivateKey {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate test key: %v", err)
	}
	return key
}

// --- Fake Memory DB Repository (Best Practice) ---
// memoryAuthRepo acts as a fast, in-memory database for our logic-layer testing.
// It avoids the brittleness of function-level mocks by accurately simulating real data state.
type memoryAuthRepo struct {
	sync.RWMutex
	users             map[int]*auth.Account
	usersByAccount    map[string]*auth.Account
	verificationCodes map[string]string
	sessions          map[string]auth.Session
	loginAudits       []auth.LoginAuditRecord
	nextID            int
	mockCache         cache.Store
}

func newMemoryAuthRepo() *memoryAuthRepo {
	return &memoryAuthRepo{
		users:             make(map[int]*auth.Account),
		usersByAccount:    make(map[string]*auth.Account),
		verificationCodes: make(map[string]string),
		sessions:          make(map[string]auth.Session),
		nextID:            1,
		mockCache:         &fakeCache{},
	}
}

// AddUser manually seeds the database for a test using this fake
func (r *memoryAuthRepo) AddUser(u *auth.Account) {
	r.Lock()
	defer r.Unlock()
	if u.UID == 0 {
		u.UID = r.nextID
		r.nextID++
	}
	r.users[u.UID] = u
	if u.Username != "" {
		r.usersByAccount[u.Username] = u
	}
	if u.Email != "" {
		r.usersByAccount[u.Email] = u
	}
	if u.Phone.Valid && u.Phone.String != "" {
		r.usersByAccount[u.Phone.String] = u
	}
}

func (r *memoryAuthRepo) GetUserByAccount(ctx context.Context, account string) (*auth.Account, error) {
	r.RLock()
	defer r.RUnlock()
	if u, exists := r.usersByAccount[account]; exists {
		return u, nil
	}
	return nil, errors.New("user not found")
}

func (r *memoryAuthRepo) GetUserByUid(ctx context.Context, uid int) (*auth.Account, error) {
	r.RLock()
	defer r.RUnlock()
	if u, exists := r.users[uid]; exists {
		return u, nil
	}
	return nil, errors.New("user not found")
}

func (r *memoryAuthRepo) RecordLoginLog(ctx context.Context, record auth.LoginAuditRecord) error {
	r.Lock()
	defer r.Unlock()
	if record.CreatedAt.IsZero() {
		record.CreatedAt = time.Now().UTC()
	}
	r.loginAudits = append(r.loginAudits, record)
	if u, ok := r.users[record.UID]; ok {
		u.LastLoginAt = &record.CreatedAt
	}
	return nil
}

func (r *memoryAuthRepo) RegisterUser(ctx context.Context, um *auth.Account) error {
	r.Lock()
	defer r.Unlock()
	// Check for unique constraint violation
	if _, exists := r.usersByAccount[um.Email]; exists && um.Email != "" {
		return errors.New("user already exists")
	}
	if _, exists := r.usersByAccount[um.Username]; exists && um.Username != "" {
		return errors.New("user already exists")
	}

	if um.UID == 0 {
		um.UID = r.nextID
		r.nextID++
	}

	r.users[um.UID] = um
	if um.Username != "" {
		r.usersByAccount[um.Username] = um
	}
	if um.Email != "" {
		r.usersByAccount[um.Email] = um
	}
	if um.Phone.Valid && um.Phone.String != "" {
		r.usersByAccount[um.Phone.String] = um
	}
	return nil
}

func (r *memoryAuthRepo) StoreVerificationCode(ctx context.Context, key, code string) error {
	r.Lock()
	defer r.Unlock()
	r.verificationCodes[key] = code
	return nil
}

func (r *memoryAuthRepo) DeleteVerificationCode(ctx context.Context, key string) error {
	r.Lock()
	defer r.Unlock()
	delete(r.verificationCodes, key)
	return nil
}

func (r *memoryAuthRepo) GetVerificationCode(ctx context.Context, key string) (string, error) {
	r.RLock()
	defer r.RUnlock()
	if code, exists := r.verificationCodes[key]; exists {
		return code, nil
	}
	return "", errors.New("code not found")
}

func (r *memoryAuthRepo) StoreUserSession(ctx context.Context, session *auth.Session, ttl time.Duration) error {
	r.Lock()
	defer r.Unlock()
	if session.CreatedAt == nil {
		t := time.Now().UTC()
		session.CreatedAt = &t
	}
	r.sessions[session.SessionID] = *session
	// Also store in mockCache for consistency with real implementation
	if r.mockCache != nil {
		b, err := json.Marshal(session)
		if err != nil {
			return err
		}
		return r.mockCache.Set(ctx, fmt.Sprintf("session:%d:%s", session.UID, session.SessionID), b, ttl)
	}
	return nil
}

func (r *memoryAuthRepo) GetUserSession(ctx context.Context, uid int, sessionID string) (*auth.Session, error) {
	r.RLock()
	defer r.RUnlock()
	if r.mockCache != nil {
		val, err := r.mockCache.Get(ctx, fmt.Sprintf("session:%d:%s", uid, sessionID))
		if err != nil {
			return nil, err
		}
		// If cache returns nil data, fall back to in-memory sessions
		if len(val) == 0 {
			if session, exists := r.sessions[sessionID]; exists {
				session.SessionID = sessionID
				return &session, nil
			}
			return nil, errors.New("not found")
		}
		var session auth.Session
		err = json.Unmarshal(val, &session)
		if err != nil {
			return nil, err
		}
		if session.CreatedAt == nil {
			t := time.Now().UTC()
			session.CreatedAt = &t
		}
		if session.SessionID == "" {
			session.SessionID = sessionID
		}
		return &session, nil
	}
	if session, exists := r.sessions[sessionID]; exists {
		if session.SessionID == "" {
			session.SessionID = sessionID
		}
		return &session, nil
	}
	return nil, nil
}

func (r *memoryAuthRepo) GetUserSessions(ctx context.Context, uid int) ([]auth.Session, error) {
	r.RLock()
	defer r.RUnlock()

	out := make([]auth.Session, 0, len(r.sessions))
	for sid, session := range r.sessions {
		if session.UID == uid {
			if session.SessionID == "" {
				session.SessionID = sid
			}
			out = append(out, session)
		}
	}
	return out, nil
}

func (r *memoryAuthRepo) DeleteUserSession(ctx context.Context, uid int, sessionID string) error {
	r.Lock()
	defer r.Unlock()
	delete(r.sessions, sessionID)
	// Also delete from mockCache for consistency with real implementation
	if r.mockCache != nil {
		return r.mockCache.Del(ctx, fmt.Sprintf("session:%d:%s", uid, sessionID))
	}
	return nil
}

func (r *memoryAuthRepo) DeleteAllUserSessions(ctx context.Context, uid int) error {
	r.Lock()
	defer r.Unlock()
	for jti, session := range r.sessions {
		if session.UID == uid {
			delete(r.sessions, jti)
		}
	}
	if r.mockCache != nil {
		return r.mockCache.DelMatch(ctx, fmt.Sprintf("session:%d:*", uid))
	}
	return nil
}

func (r *memoryAuthRepo) UpdateUserPassword(ctx context.Context, uid int, password string) error {
	r.Lock()
	defer r.Unlock()
	if user, ok := r.users[uid]; ok {
		user.Password = password
		return nil
	}
	return errors.New("user not found")
}

func (r *memoryAuthRepo) StorePasswordResetToken(ctx context.Context, token, email string) error {
	r.Lock()
	defer r.Unlock()
	r.verificationCodes["reset_token:"+token] = email
	return nil
}

func (r *memoryAuthRepo) GetPasswordResetEmail(ctx context.Context, token string) (string, error) {
	r.RLock()
	defer r.RUnlock()
	if email, exists := r.verificationCodes["reset_token:"+token]; exists {
		return email, nil
	}
	return "", errors.New("token not found")
}

func (r *memoryAuthRepo) DeletePasswordResetToken(ctx context.Context, token string) error {
	r.Lock()
	defer r.Unlock()
	delete(r.verificationCodes, "reset_token:"+token)
	return nil
}

func (r *memoryAuthRepo) Cache() cache.Store {
	return r.mockCache
}

func (r *memoryAuthRepo) ListRecentLogins(ctx context.Context, uid int, limit int) ([]auth.LoginAuditRecord, error) {
	r.RLock()
	defer r.RUnlock()
	out := make([]auth.LoginAuditRecord, 0)
	for _, record := range r.loginAudits {
		if record.UID == uid {
			out = append(out, record)
		}
	}
	if limit > 0 && len(out) > limit {
		out = out[len(out)-limit:]
	}
	return out, nil
}

// fakeCache is a simple no-op cache for testing
type fakeCache struct{}

func (m *fakeCache) Get(ctx context.Context, key string) ([]byte, error) { return nil, nil }
func (m *fakeCache) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return nil
}
func (m *fakeCache) Del(ctx context.Context, key string) error          { return nil }
func (m *fakeCache) DelMatch(ctx context.Context, pattern string) error { return nil }
func (m *fakeCache) GetMatch(ctx context.Context, pattern string) (map[string][]byte, error) {
	return make(map[string][]byte), nil
}
func (m *fakeCache) Flush(ctx context.Context) error { return nil }

var _ dao.AuthRepository = (*memoryAuthRepo)(nil)

type memoryRBACReader struct{}

func (r *memoryRBACReader) LoadUserRoles(ctx context.Context, uid int) ([]rbac.Role, error) {
	return []rbac.Role{{ID: 1, Code: "User"}}, nil
}

func (r *memoryRBACReader) LoadUserPermissions(ctx context.Context, uid int) ([]rbac.Permission, error) {
	return []rbac.Permission{}, nil
}

func (r *memoryRBACReader) LoadRolePermissions(ctx context.Context, role string) ([]rbac.Permission, error) {
	return []rbac.Permission{}, nil
}

var _ rbac.Reader = (*memoryRBACReader)(nil)
