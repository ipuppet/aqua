package token

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type revokeCacheStub struct {
	data   map[string][]byte
	setErr error
}

func (s *revokeCacheStub) Get(ctx context.Context, key string) ([]byte, error) {
	if v, ok := s.data[key]; ok {
		return v, nil
	}
	return nil, redis.Nil
}

func (s *revokeCacheStub) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	if s.setErr != nil {
		return s.setErr
	}
	if s.data == nil {
		s.data = map[string][]byte{}
	}
	s.data[key] = []byte(value.(string))
	return nil
}

func (s *revokeCacheStub) Del(ctx context.Context, key string) error {
	delete(s.data, key)
	return nil
}

func (s *revokeCacheStub) DelMatch(ctx context.Context, pattern string) error { return nil }
func (s *revokeCacheStub) GetMatch(ctx context.Context, pattern string) (map[string][]byte, error) {
	return make(map[string][]byte), nil
}
func (s *revokeCacheStub) Flush(ctx context.Context) error { return nil }

func TestRevokeUserRefreshTokens_WritesMarker(t *testing.T) {
	c := &revokeCacheStub{data: map[string][]byte{}}
	err := RevokeUserRefreshTokens(context.Background(), c, 42)
	require.NoError(t, err)

	v, ok := c.data[refreshTokenRevokeAfterUserKey(42)]
	require.True(t, ok)
	ts, parseErr := strconv.ParseInt(string(v), 10, 64)
	require.NoError(t, parseErr)
	assert.NotZero(t, ts)
}

func TestRevokeUserRefreshTokens_PropagatesSetError(t *testing.T) {
	expected := errors.New("set failed")
	c := &revokeCacheStub{setErr: expected}
	err := RevokeUserRefreshTokens(context.Background(), c, 42)
	assert.ErrorIs(t, err, expected)
}

func TestRevokeAllRefreshTokens_PropagatesSetError(t *testing.T) {
	expected := errors.New("set failed")
	c := &revokeCacheStub{setErr: expected}
	err := RevokeAllRefreshTokens(context.Background(), c)
	assert.ErrorIs(t, err, expected)
}

func TestIsRefreshTokenRevoked_UsesUserOrGlobalMaxMarker(t *testing.T) {
	c := &revokeCacheStub{data: map[string][]byte{
		refreshTokenRevokeAfterGlobalKey(): []byte("100"),
		refreshTokenRevokeAfterUserKey(42): []byte("150"),
	}}

	revoked, err := IsRefreshTokenRevoked(context.Background(), c, 42, 149)
	require.NoError(t, err)
	assert.True(t, revoked)

	notRevoked, err := IsRefreshTokenRevoked(context.Background(), c, 42, 151)
	require.NoError(t, err)
	assert.False(t, notRevoked)
}
