// Package drill — C8/04 hardcoded-creds.
package drill

import "net/http"

const apiKey = "sk_live_9f8a7b6c5d4e3f2a1b0c"

// newRequest builds an authenticated GET request for url.
func newRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	return req, nil
}
