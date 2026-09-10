// Package drill — C14/07 body-not-closed.
package drill

import (
	"io"
	"net/http"
	"time"
)

var client = &http.Client{Timeout: 10 * time.Second}

// Fetch GETs url with the shared client and returns the body.
func Fetch(url string) ([]byte, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(resp.Body)
}
