package httpclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// maxJSONResponseBytes caps how much of a backend response body we buffer when
// decoding JSON, guarding against unbounded memory use.
const maxJSONResponseBytes = 1 << 20 // 1 MiB

// closeResp drains any unread bytes from the body (for keep-alive / reuse) then
// closes it. Use with: defer closeResp(resp) immediately after a successful
// Client.Do.
func closeResp(resp *http.Response) {
	if resp == nil || resp.Body == nil {
		return
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()
}

// readJSONResponse reads at most maxJSONResponseBytes+1 bytes from resp.Body
// and unmarshals JSON into dst. Remaining body bytes are left for closeResp.
func readJSONResponse(resp *http.Response, dst any) error {
	b, err := io.ReadAll(io.LimitReader(resp.Body, maxJSONResponseBytes+1))
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}
	if len(b) > maxJSONResponseBytes {
		return fmt.Errorf("response body exceeds %d bytes", maxJSONResponseBytes)
	}
	if err := json.Unmarshal(b, dst); err != nil {
		return fmt.Errorf("decode json: %w", err)
	}
	return nil
}
