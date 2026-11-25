// Package client provides HTTP clients with IPv4 dialing and OpenTelemetry instrumentation.
package client

import (
	"context"
	"net"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// IPv4Transport returns a copy of the default transport hardcoded to prefer strictly IPv4 resolutions.
// 在某些纯 IPv4 的容器网络编排环境下强制切断双栈解析，能消除因 IPv6 DNS 解析无端回退引发的网络超时问题。
func IPv4Transport() *http.Transport {
	base := http.DefaultTransport.(*http.Transport).Clone()
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	base.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		return dialer.DialContext(ctx, "tcp4", addr)
	}
	return base
}

// NewHTTPClient constructs an [http.Client] bundled with IPv4 dialing and tracing instrumentation.
// 直接在 Transport 层包裹一层 otelhttp 插件，可在无需修改上层业务调用的前提下自动分发和记录 OpenTelemetry 链路 Span 信息。
func NewHTTPClient(timeout time.Duration) *http.Client {
	baseTransport := IPv4Transport()
	instrumentedTransport := otelhttp.NewTransport(baseTransport)
	return &http.Client{
		Timeout:   timeout,
		Transport: instrumentedTransport,
	}
}
