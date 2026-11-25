package authctx

import (
	"context"
	"net"
	"net/http"
	"strings"

	"google.golang.org/grpc/metadata"
)

const metadataClientIP = "x-client-ip"

// WithRequestClientIP copies the remote address into gRPC incoming metadata.
func WithRequestClientIP(_ context.Context, r *http.Request) metadata.MD {
	ip := requestClientIP(r)
	if ip == "" {
		return nil
	}
	return metadata.Pairs(metadataClientIP, ip)
}

// ClientIP returns the client address previously stored by WithRequestClientIP.
func ClientIP(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	vals := md.Get(metadataClientIP)
	if len(vals) == 0 {
		return ""
	}
	return vals[0]
}

func requestClientIP(r *http.Request) string {
	if r == nil {
		return ""
	}
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return strings.TrimSpace(strings.Split(xff, ",")[0])
	}
	if realIP := strings.TrimSpace(r.Header.Get("X-Real-IP")); realIP != "" {
		return realIP
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
