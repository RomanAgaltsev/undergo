// Package drill — C14/02 status-unchecked.
package drill

import (
	"io"
	"net/http"
	"time"
)

var client = &http.Client{Timeout: 10 * time.Second}

// GetJSON fetches url with a shared, timed client and returns the response body.
func GetJSON(url string) ([]byte, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
