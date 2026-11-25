package authctx

import (
	"context"
	"net/http"
	"testing"

	"google.golang.org/grpc/metadata"
)

func TestClientIPFromForwardedHeader(t *testing.T) {
	t.Parallel()

	req, err := http.NewRequest(http.MethodGet, "http://example.com/api/articles/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("X-Forwarded-For", "203.0.113.9, 10.0.0.1")
	req.RemoteAddr = "192.0.2.1:443"

	md := WithRequestClientIP(context.Background(), req)
	ctx := metadata.NewIncomingContext(context.Background(), md)
	if got := ClientIP(ctx); got != "203.0.113.9" {
		t.Fatalf("ClientIP = %q, want 203.0.113.9", got)
	}
}

func TestClientIPFallsBackToRemoteAddr(t *testing.T) {
	t.Parallel()

	req, err := http.NewRequest(http.MethodGet, "http://example.com/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.RemoteAddr = "198.51.100.10:5150"

	md := WithRequestClientIP(context.Background(), req)
	ctx := metadata.NewIncomingContext(context.Background(), md)
	if got := ClientIP(ctx); got != "198.51.100.10" {
		t.Fatalf("ClientIP = %q, want 198.51.100.10", got)
	}
}
