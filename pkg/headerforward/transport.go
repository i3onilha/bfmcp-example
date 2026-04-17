// Package headerforward provides an HTTP transport that injects headers
// from a context into outgoing HTTP requests.
package headerforward

import (
	"net/http"
)

// ContextKey is the context key used to propagate headers.
type ContextKey struct{}

// Transport is an http.RoundTripper that injects headers
// from the context into the outgoing request.
type Transport struct{}

// RoundTrip executes a single HTTP transaction, adding headers from the context.
func (h *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	newReq := req.Clone(req.Context())

	if headers, ok := req.Context().Value(ContextKey{}).(http.Header); ok {
		for key, values := range headers {
			newReq.Header.Del(key)
			for _, value := range values {
				newReq.Header.Add(key, value)
			}
		}
	}

	return http.DefaultTransport.RoundTrip(newReq)
}

// NewClient creates a new http.Client that uses the header forwarding transport.
func NewClient() *http.Client {
	return &http.Client{
		Transport: &Transport{},
	}
}
