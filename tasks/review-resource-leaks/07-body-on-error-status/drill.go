// Package drill — C5/07 body-on-error-status.
package drill

import (
	"fmt"
	"io"
	"net/http"
)

// Get fetches url and returns the body, erroring on a non-200 status.
func Get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
