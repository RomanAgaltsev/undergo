// Package drill — C6/08 client-api.
package drill

import (
	"sync"
	"time"
)

// Doer performs requests.
type Doer interface {
	Do(req string) string
}

// Client is an API client.
type Client struct {
	Mu      sync.Mutex
	timeout time.Duration
}

// New returns a ready client.
func New() Doer {
	return &Client{timeout: 5 * time.Second}
}

// Do performs a request.
func (c *Client) Do(req string) string {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	return req
}

// Timeout reports the configured timeout.
func (c Client) Timeout() time.Duration {
	return c.timeout
}
