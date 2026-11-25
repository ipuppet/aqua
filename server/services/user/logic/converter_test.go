package logic

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	user "aqua/services/user/domain"
)

func TestToUserModelResponseKeepsTimestamps(t *testing.T) {
	registered := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
	lastLogin := time.Date(2026, 8, 24, 10, 0, 0, 0, time.UTC)
	got := ToUserModelResponse(&user.UserModel{
		UID:         1,
		Username:    "alice",
		RegisterAt:  &registered,
		LastLoginAt: &lastLogin,
	})
	require.NotNil(t, got)
	require.NotNil(t, got.RegisterAt)
	require.NotNil(t, got.LastLoginAt)
	assert.True(t, registered.Equal(*got.RegisterAt))
	assert.True(t, lastLogin.Equal(*got.LastLoginAt))
}
