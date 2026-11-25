package blog

import (
	"time"

	"aqua/pkg/types"
)

// AccessResponse represents the API response structure for an access log record.
type AccessResponse struct {
	ID             int        `json:"id"`
	Aid            int        `json:"aid"`
	IPv4           string     `json:"ipv4"`
	IPv6           string     `json:"ipv6"`
	Timestamp      *time.Time `json:"timestamp"`
	TimestampEpoch string     `json:"timestamp_epoch"`
	Title          string     `json:"title"`
}

// Access represents a database record documenting a single client visit.
type Access struct {
	ID             int              `db:"id" json:"id"`
	Aid            int              `db:"aid" json:"aid"`
	IPv4           types.NullString `db:"ipv4" json:"ipv4"`
	IPv6           types.NullString `db:"ipv6" json:"ipv6"`
	Timestamp      *time.Time       `db:"timestamp" json:"timestamp"`
	TimestampEpoch string           `db:"timestamp_epoch" json:"timestamp_epoch"`
	Title          types.NullString `db:"title" json:"title"`
}
