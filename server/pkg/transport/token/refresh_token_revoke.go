package token

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/ipuppet/aqua/server/pkg/infra/cache"

	"github.com/redis/go-redis/v9"
)

const (
	refreshTokenRevokeAfterUserKeyPre     = "refresh_token_revoke_after:user:"
	refreshTokenRevokeAfterGlobalKeyValue = "refresh_token_revoke_after:global"
	// refreshTokenRevokeMarkerTTL is the TTL for per-user and global revoke-after markers (stored as Unix seconds).
	// Must be strictly greater than the refresh token lifetime configured in services/passport/logic/passport.go
	// (currently refreshTokenLife = 7 * 24h). Update both values together when changing token lifetime.
	refreshTokenRevokeMarkerTTL = 8 * 24 * time.Hour
)

// refreshTokenRevokeAfterUserKey returns the cache key for a per-user revoke-after marker.
func refreshTokenRevokeAfterUserKey(uid int) string {
	return refreshTokenRevokeAfterUserKeyPre + strconv.Itoa(uid)
}

// refreshTokenRevokeAfterGlobalKey returns the cache key for the global revoke-after marker.
func refreshTokenRevokeAfterGlobalKey() string {
	return refreshTokenRevokeAfterGlobalKeyValue
}

// RevokeUserRefreshTokens writes a per-user revoke-after marker. Any refresh token
// issued before this call will be rejected on next use.
func RevokeUserRefreshTokens(ctx context.Context, c cache.Store, uid int) error {
	if c == nil || uid <= 0 {
		return nil
	}
	v := strconv.FormatInt(time.Now().Unix(), 10)
	return c.Set(ctx, refreshTokenRevokeAfterUserKey(uid), v, refreshTokenRevokeMarkerTTL)
}

// RevokeAllRefreshTokens writes a global revoke-after marker. Any refresh token
// issued before this call will be rejected on next use.
func RevokeAllRefreshTokens(ctx context.Context, c cache.Store) error {
	if c == nil {
		return nil
	}
	v := strconv.FormatInt(time.Now().Unix(), 10)
	return c.Set(ctx, refreshTokenRevokeAfterGlobalKey(), v, refreshTokenRevokeMarkerTTL)
}

// IsRefreshTokenRevoked reports whether the refresh token issued at issuedAt (Unix seconds)
// for uid has been superseded by a per-user or global revoke-after marker.
func IsRefreshTokenRevoked(ctx context.Context, c cache.Store, uid int, issuedAt int64) (bool, error) {
	if c == nil {
		return false, nil
	}
	maxRevokeAfter, err := readMaxRevokeAfter(ctx, c,
		refreshTokenRevokeAfterGlobalKey(),
		refreshTokenRevokeAfterUserKey(uid),
	)
	if err != nil {
		return false, err
	}
	if maxRevokeAfter == 0 {
		return false, nil
	}
	return issuedAt <= maxRevokeAfter, nil
}

func readMaxRevokeAfter(ctx context.Context, c cache.Store, keys ...string) (int64, error) {
	var maxTs int64
	for _, key := range keys {
		b, err := c.Get(ctx, key)
		if err != nil {
			if errors.Is(err, redis.Nil) {
				continue
			}
			return 0, err
		}
		ts, err := strconv.ParseInt(string(b), 10, 64)
		if err != nil {
			continue
		}
		if ts > maxTs {
			maxTs = ts
		}
	}
	return maxTs, nil
}
