// Package drill — C14/01 no-timeout.
package drill

import (
	"io"
	"net/http"
)

// Fetch GETs url and returns the body bytes.
func Fetch(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
