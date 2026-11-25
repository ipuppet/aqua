package grpcx

import (
	"context"
	"errors"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ipuppet/aqua/server/pkg/apperr"
)

// MapError converts apperr (and similar) values into gRPC status errors.
func MapError(err error) error {
	if err == nil {
		return nil
	}

	var mapped apperr.GRPCMappableError
	if errors.As(err, &mapped) {
		return status.Error(mapped.GRPCCode(), mapped.ClientMessage())
	}

	return err
}

// UnaryServerErrorMapping maps returned application errors into gRPC status codes.
func UnaryServerErrorMapping() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			return nil, MapError(err)
		}
		return resp, nil
	}
}

var grpcToHTTPStatus = map[codes.Code]int{
	codes.InvalidArgument:    http.StatusBadRequest,
	codes.Unauthenticated:    http.StatusUnauthorized,
	codes.PermissionDenied:   http.StatusForbidden,
	codes.NotFound:           http.StatusNotFound,
	codes.AlreadyExists:      http.StatusConflict,
	codes.FailedPrecondition: http.StatusBadRequest,
	codes.Aborted:            http.StatusConflict,
	codes.OutOfRange:         http.StatusBadRequest,
	codes.ResourceExhausted:  http.StatusTooManyRequests,
	codes.Canceled:           499,
	codes.DeadlineExceeded:   http.StatusGatewayTimeout,
	codes.Unimplemented:      http.StatusNotImplemented,
	codes.Unavailable:        http.StatusServiceUnavailable,
	codes.DataLoss:           http.StatusInternalServerError,
	codes.Unknown:            http.StatusInternalServerError,
	codes.Internal:           http.StatusInternalServerError,
}

// HTTPStatusFromCode converts a gRPC code into the corresponding HTTP status.
func HTTPStatusFromCode(code codes.Code) int {
	if httpCode, exists := grpcToHTTPStatus[code]; exists {
		return httpCode
	}
	return http.StatusInternalServerError
}
