// Package drill — C1/07 once-copied.
package drill

import "sync"

// Resource lazily initializes its connection exactly once.
type Resource struct {
	once sync.Once
	conn string
}

// Get returns the resource's connection, initializing it on first use.
func (r Resource) Get() string {
	r.once.Do(func() {
		r.conn = "connected"
	})
	return r.conn
}
