package logic

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"aqua/pkg/types"
	blog "aqua/services/blog/domain"
)

func TestToAccessResponseKeepsTimestamp(t *testing.T) {
	ts := time.Date(2026, 8, 24, 14, 30, 0, 0, time.UTC)
	got := ToAccessResponse(&blog.Access{
		ID:             9,
		Aid:            3,
		IPv4:           types.NullString{NullString: sql.NullString{String: "127.0.0.1", Valid: true}},
		Timestamp:      &ts,
		TimestampEpoch: "1756045800",
		Title:          types.NullString{NullString: sql.NullString{String: "hello", Valid: true}},
	})
	require.NotNil(t, got)
	require.NotNil(t, got.Timestamp)
	assert.True(t, ts.Equal(*got.Timestamp))
	assert.Equal(t, "hello", got.Title)
	assert.Equal(t, "127.0.0.1", got.IPv4)
}
