package dao

import (
	"context"
	"errors"
	"testing"
	"time"

	auth "github.com/ipuppet/aqua/server/services/passport/domain"

	"github.com/stretchr/testify/require"
)

type sessionTestCache struct {
	data map[string][]byte
}

func (c *sessionTestCache) Get(ctx context.Context, key string) ([]byte, error) {
	if v, ok := c.data[key]; ok {
		return v, nil
	}
	return nil, errors.New("not found")
}

func (c *sessionTestCache) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	if c.data == nil {
		c.data = map[string][]byte{}
	}
	switch v := value.(type) {
	case []byte:
		c.data[key] = v
	case string:
		c.data[key] = []byte(v)
	default:
		return errors.New("unsupported cache value type")
	}
	return nil
}

func (c *sessionTestCache) Del(ctx context.Context, key string) error          { return nil }
func (c *sessionTestCache) DelMatch(ctx context.Context, pattern string) error { return nil }
func (c *sessionTestCache) GetMatch(ctx context.Context, pattern string) (map[string][]byte, error) {
	return map[string][]byte{}, nil
}
func (c *sessionTestCache) Flush(ctx context.Context) error { return nil }

func TestGetUserSession_UnmarshalValue(t *testing.T) {
	ctx := context.Background()
	createdAt := time.Now().UTC().Round(time.Second)
	repo := &repository{cache: &sessionTestCache{data: map[string][]byte{}}}

	s := &auth.Session{
		SessionID:  "sid-1",
		UID:        42,
		DeviceInfo: "ios",
		CreatedAt:  &createdAt,
	}
	require.NoError(t, repo.StoreUserSession(ctx, s, time.Minute))

	got, err := repo.GetUserSession(ctx, 42, "sid-1")
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, s.SessionID, got.SessionID)
	require.Equal(t, s.UID, got.UID)
	require.Equal(t, s.DeviceInfo, got.DeviceInfo)
	require.NotNil(t, got.CreatedAt)
}
