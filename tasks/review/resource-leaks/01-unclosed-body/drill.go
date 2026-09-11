// Package drill — C5/01 unclosed-body.
package drill

import (
	"io"
	"net/http"
)

// Fetch performs a GET and returns the response body bytes.
func Fetch(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(resp.Body)
}

// FetchString performs a GET and returns the body as a string.
func FetchString(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
