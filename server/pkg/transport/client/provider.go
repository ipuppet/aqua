package client

import (
	"net/http"
	"time"

	"github.com/google/wire"
)

// DefaultHTTPClient returns an HTTP client with default timeout (for wire).
func DefaultHTTPClient() *http.Client {
	return NewHTTPClient(15 * time.Second)
}

// Set is client providers.
var Set = wire.NewSet(DefaultHTTPClient)
