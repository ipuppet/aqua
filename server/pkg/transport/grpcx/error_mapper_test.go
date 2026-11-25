package grpcx

import (
	"errors"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"aqua/pkg/apperr"
)

func TestMapError_MapsDomainErrorToGRPCStatus(t *testing.T) {
	t.Parallel()

	err := MapError(apperr.Forbidden("forbidden"))
	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("expected grpc status error")
	}

	if st.Code() != codes.PermissionDenied {
		t.Fatalf("unexpected grpc code: got %s, want %s", st.Code(), codes.PermissionDenied)
	}
}

func TestMapError_PassThroughUnknownError(t *testing.T) {
	t.Parallel()

	raw := errors.New("boom")
	if got := MapError(raw); got != raw {
		t.Fatalf("expected same error instance for unknown error")
	}
}
