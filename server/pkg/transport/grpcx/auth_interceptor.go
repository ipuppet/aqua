package grpcx

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const InternalTokenMetadataKey = "x-internal-token"

// UnaryServerInternalAuth validates the internal service token for unary RPC calls.
// The expectedTokenProvider is evaluated per request so runtime config changes can take effect.
func UnaryServerInternalAuth(expectedTokenProvider func() string) grpc.UnaryServerInterceptor {
	if expectedTokenProvider == nil {
		return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
			return nil, status.Errorf(codes.Internal, "internal token provider is not configured")
		}
	}
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		expectedToken := strings.TrimSpace(expectedTokenProvider())
		if expectedToken == "" {
			return nil, status.Errorf(codes.Internal, "internal token is not configured")
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "metadata is required")
		}

		tokens := md.Get(InternalTokenMetadataKey)
		if len(tokens) == 0 || strings.TrimSpace(tokens[0]) == "" {
			return nil, status.Errorf(codes.Unauthenticated, "internal token is required")
		}

		if tokens[0] != expectedToken {
			return nil, status.Errorf(codes.PermissionDenied, "invalid internal token")
		}

		return handler(ctx, req)
	}
}
