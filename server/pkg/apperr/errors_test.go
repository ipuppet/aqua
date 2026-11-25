package apperr

import (
	"database/sql"
	"errors"
	"testing"
)

func TestFromDAO(t *testing.T) {
	t.Parallel()

	if got := FromDAO(sql.ErrNoRows); got == nil {
		t.Fatal("expected NotFound, got nil")
	} else {
		var mapped *Error
		if !errors.As(got, &mapped) {
			t.Fatalf("expected *Error, got %T", got)
		}
		if mapped.HTTPStatus() != 404 {
			t.Fatalf("HTTP status: got %d, want 404", mapped.HTTPStatus())
		}
	}

	orig := errors.New("disk full")
	if got := FromDAO(orig); !errors.Is(got, orig) {
		t.Fatalf("non-no-rows error should pass through, got %v", got)
	}
}
