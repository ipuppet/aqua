// Package gateway provides grpc-gateway ServeMux options.
package gateway

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/encoding/protojson"

	"aqua/pkg/transport/authctx"
)

// NewServeMux creates a [runtime.ServeMux] pre-configured with standard options.
//
// Extra opts are appended after the default options, so callers can still
// add WithMetadata, WithIncomingHeaderMatcher, etc.
func NewServeMux(opts ...runtime.ServeMuxOption) *runtime.ServeMux {
	defaults := []runtime.ServeMuxOption{
		runtime.WithMetadata(authctx.WithRequestClientIP),
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:   false,
				EmitUnpopulated: false,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	}
	return runtime.NewServeMux(append(defaults, opts...)...)
}
