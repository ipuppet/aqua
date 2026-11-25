package response

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/ipuppet/aqua/server/pkg/apperr"
)

func TestFail_MapsDomainErrors(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name     string
		err      error
		wantCode int
	}{
		{
			name:     "forbidden",
			err:      apperr.Forbidden("forbidden"),
			wantCode: http.StatusForbidden,
		},
		{
			name:     "unauthorized",
			err:      apperr.Unauthorized(""),
			wantCode: http.StatusUnauthorized,
		},
		{
			name:     "bad request",
			err:      apperr.BadRequest("bad request"),
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "not found",
			err:      apperr.NotFound("resource missing"),
			wantCode: http.StatusNotFound,
		},
		{
			name:     "unavailable",
			err:      apperr.ServiceUnavailable("unavailable"),
			wantCode: http.StatusServiceUnavailable,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)

			Fail(ctx, tc.err)

			if recorder.Code != tc.wantCode {
				t.Fatalf("unexpected status code: got %d, want %d", recorder.Code, tc.wantCode)
			}
		})
	}
}

func TestFail_UnknownErrorReturnsInternalServerError(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	Fail(ctx, errors.New("boom"))

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("unexpected status code: got %d, want %d", recorder.Code, http.StatusInternalServerError)
	}
}
