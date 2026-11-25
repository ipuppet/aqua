// Package apperr is the shared application error vocabulary.
// A single error carries the client message plus both HTTP status and gRPC code
// so Gin (response.Fail) and gRPC (grpcx.MapError) map from the same value.
package apperr

import (
	"database/sql"
	"errors"

	"google.golang.org/grpc/codes"
)

// Error is an application error with HTTP and gRPC mappings.
type Error struct {
	message    string
	httpStatus int
	grpcCode   codes.Code
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	return e.message
}

// HTTPStatus maps this error to an HTTP response status.
func (e *Error) HTTPStatus() int {
	if e == nil {
		return 500
	}
	return e.httpStatus
}

// ClientMessage returns client-safe message.
func (e *Error) ClientMessage() string {
	if e == nil {
		return "internal server error"
	}
	return e.message
}

// GRPCCode maps this error to a gRPC status code.
func (e *Error) GRPCCode() codes.Code {
	if e == nil {
		return codes.Internal
	}
	return e.grpcCode
}

func newError(message string, httpStatus int, grpcCode codes.Code) error {
	return &Error{
		message:    message,
		httpStatus: httpStatus,
		grpcCode:   grpcCode,
	}
}

func BadRequest(message string) error {
	return newError(message, 400, codes.InvalidArgument)
}

func Unauthorized(message string) error {
	if message == "" {
		message = "unauthorized access"
	}
	return newError(message, 401, codes.Unauthenticated)
}

func Forbidden(message string) error {
	return newError(message, 403, codes.PermissionDenied)
}

func NotFound(message string) error {
	return newError(message, 404, codes.NotFound)
}

func ServiceUnavailable(message string) error {
	return newError(message, 503, codes.Unavailable)
}

func Internal(message string) error {
	if message == "" {
		message = "internal server error"
	}
	return newError(message, 500, codes.Internal)
}

// FromDAO maps a storage miss to NotFound and otherwise returns err unchanged.
func FromDAO(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return NotFound("record not found")
	}
	return err
}
