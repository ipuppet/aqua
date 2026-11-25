package dao

import (
	"context"
	"time"
)

func passwordResetTokenKey(token string) string {
	return "password_reset_token:" + token
}

// StorePasswordResetToken saves the password reset token with the associated email in Redis with a 15-minute TTL.
func (r *repository) StorePasswordResetToken(ctx context.Context, token, email string) error {
	return r.cache.Set(ctx, passwordResetTokenKey(token), email, 15*time.Minute)
}

func (r *repository) GetPasswordResetEmail(ctx context.Context, token string) (string, error) {
	val, err := r.cache.Get(ctx, passwordResetTokenKey(token))
	if err != nil {
		return "", err
	}
	return string(val), nil
}

func (r *repository) DeletePasswordResetToken(ctx context.Context, token string) error {
	return r.cache.Del(ctx, passwordResetTokenKey(token))
}
