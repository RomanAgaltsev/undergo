// Package stamp calls a pointer method on every element of a slice of values.
//
// The caller holds a []Event. Event does not implement Stampable — only *Event
// does — and the elements must be modified in place, not copied. Getting that
// past the type checker without copying is the task.
package stamp

import "time"

// Stampable is implemented on the POINTER receiver, which is the whole problem.
type Stampable interface {
	SetStamp(at time.Time)
}

// Event is the value type the tests use. Note the receiver below.
type Event struct {
	Name    string
	Stamped time.Time
}

// SetStamp records when the event was stamped.
func (e *Event) SetStamp(at time.Time) { e.Stamped = at }

// StampAll stamps every item in place and returns the same slice.
//
// PT is constrained twice over: it must be the pointer to T, and it must satisfy
// Stampable. Both halves are load-bearing — the first lets you take the address
// of an element and convert it, the second lets you call the method.
//
// The contract the frozen tests enforce:
//
//   - every element of the ORIGINAL slice is stamped, not a copy of it;
//   - the returned slice shares its backing array with the argument;
//   - a nil or empty slice returns without panicking;
//   - the whole walk allocates nothing.
func StampAll[T any, PT interface {
	*T
	Stampable
}](items []T, at time.Time) []T {
	panic("undergo: implement StampAll")
}
