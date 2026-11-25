package convert

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCloneTime(t *testing.T) {
	assert.Nil(t, CloneTime(nil))

	src := time.Date(2026, 3, 4, 5, 6, 7, 8, time.UTC)
	got := CloneTime(&src)
	require.NotNil(t, got)
	assert.True(t, src.Equal(*got))
	assert.NotSame(t, &src, got)
}
