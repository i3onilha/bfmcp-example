package headerforward_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"bfm-example/pkg/headerforward"
)

func TestTransport_replacesForwardedHeaders(t *testing.T) {
	t.Parallel()
	var gotTenant string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotTenant = r.Header.Get("X-Tenant-Id")
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)

	ctx := context.WithValue(context.Background(), headerforward.ContextKey{}, headerforward.FilterHeaders(http.Header{
		"X-Tenant-Id": []string{"from-context"},
	}))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, srv.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("X-Tenant-Id", "from-original-request")

	client := headerforward.NewClient()
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()

	if gotTenant != "from-context" {
		t.Fatalf("X-Tenant-Id = %q, want from-context (forwarded values must replace, not append)", gotTenant)
	}
}
