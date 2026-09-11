// Package drill — C6/01 value-receiver.
package drill

// Counter counts events.
type Counter struct {
	n int
}

// Inc increments the counter.
func (c Counter) Inc() {
	c.n++
}

// Value returns the current count.
func (c Counter) Value() int {
	return c.n
}
