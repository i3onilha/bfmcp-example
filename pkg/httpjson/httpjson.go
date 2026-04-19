// Package httpjson provides bounded JSON decoding from HTTP responses.
package httpjson

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// MaxJSONBodyBytes caps how much of a response body we buffer when decoding
// JSON, guarding against unbounded memory use.
const MaxJSONBodyBytes = 1 << 20 // 1 MiB

// DrainAndClose drains any unread bytes from the body (for keep-alive / reuse)
// then closes it. Defer immediately after a successful Client.Do when the
// response body must always be released.
func DrainAndClose(resp *http.Response) {
	if resp == nil || resp.Body == nil {
		return
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()
}

// DecodeJSON reads at most MaxJSONBodyBytes+1 bytes from resp.Body and
// unmarshals JSON into dst. Remaining body bytes are left for DrainAndClose.
func DecodeJSON(resp *http.Response, dst any) error {
	b, err := io.ReadAll(io.LimitReader(resp.Body, MaxJSONBodyBytes+1))
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}
	if len(b) > MaxJSONBodyBytes {
		return fmt.Errorf("response body exceeds %d bytes", MaxJSONBodyBytes)
	}
	if err := json.Unmarshal(b, dst); err != nil {
		return fmt.Errorf("decode json: %w", err)
	}
	return nil
}
