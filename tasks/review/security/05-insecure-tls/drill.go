// Package drill — C8/05 insecure-tls.
package drill

import (
	"crypto/tls"
	"net/http"
	"time"
)

// newClient returns an HTTP client for talking to the upstream API.
func newClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}
