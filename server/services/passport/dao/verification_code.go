package dao

import (
	"context"
	"time"
)

// StoreVerificationCode persists the verification code with a five minute TTL.
func (r *repository) StoreVerificationCode(ctx context.Context, key, code string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return r.cache.Set(ctx, key, code, 5*time.Minute)
}

// DeleteVerificationCode removes the cached verification code.
func (r *repository) DeleteVerificationCode(ctx context.Context, key string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return r.cache.Del(ctx, key)
}

// GetVerificationCode fetches a cached verification code.
func (r *repository) GetVerificationCode(ctx context.Context, key string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	val, err := r.cache.Get(ctx, key)
	if err != nil {
		return "", err
	}
	return string(val), nil
}
