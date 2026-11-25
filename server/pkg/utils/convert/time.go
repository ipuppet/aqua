package convert

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToPBTime converts a *time.Time to a *timestamppb.Timestamp.
func ToPBTime(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}

// CloneTime copies t. goverter cannot map time.Time under ignoreUnexported
// because time.Time fields are unexported; without this extend it emits a zero time.
func CloneTime(t *time.Time) *time.Time {
	if t == nil {
		return nil
	}
	copied := *t
	return &copied
}
