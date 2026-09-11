// Package drill — C14/08 drain-before-close.
package drill

import (
	"encoding/json"
	"net/http"
	"time"
)

var client = &http.Client{Timeout: 10 * time.Second}

// Ping is the decoded health response.
type Ping struct {
	OK bool `json:"ok"`
}

// Health fetches the health endpoint and decodes the first JSON object.
func Health(url string) (Ping, error) {
	resp, err := client.Get(url)
	if err != nil {
		return Ping{}, err
	}

	var p Ping
	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		resp.Body.Close()
		return Ping{}, err
	}
	resp.Body.Close()

	return p, nil
}
