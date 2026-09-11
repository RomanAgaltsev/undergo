// Package drill — C14/03 client-per-request.
package drill

import (
	"io"
	"net/http"
	"time"
)

// Get fetches url and returns the body.
func Get(url string) ([]byte, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
