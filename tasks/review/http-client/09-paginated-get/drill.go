// Package drill — C14/09 paginated-get.
package drill

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// FetchAll fetches `pages` pages from base and returns their raw bodies.
func FetchAll(ctx context.Context, base string, pages int) ([][]byte, error) {
	out := make([][]byte, 0, pages)

	for i := 0; i < pages; i++ {
		client := &http.Client{Timeout: 10 * time.Second}

		url := fmt.Sprintf("%s?page=%d", base, i)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		out = append(out, body)
	}

	return out, nil
}
