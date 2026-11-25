package apperr

import "google.golang.org/grpc/codes"

// HTTPMappableError maps an application error to HTTP.
type HTTPMappableError interface {
	error
	HTTPStatus() int
	ClientMessage() string
}

// GRPCMappableError maps an application error to gRPC.
type GRPCMappableError interface {
	error
	GRPCCode() codes.Code
	ClientMessage() string
}
