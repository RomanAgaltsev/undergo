// Package drill — C14/04 no-context.
package drill

import (
	"context"
	"io"
	"net/http"
	"time"
)

var client = &http.Client{Timeout: 30 * time.Second}

// Fetch fetches url on behalf of the caller, who supplies a context for
// cancellation and deadlines.
func Fetch(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
