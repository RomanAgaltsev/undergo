// Package boundedqueue implements a bounded blocking FIFO on sync.Cond.
package boundedqueue

import "errors"

// ErrClosed is returned by Push and Pop once the queue is closed and can no
// longer satisfy the call.
var ErrClosed = errors.New("queue is closed")

// Queue is a bounded FIFO. Push blocks while the queue is full, Pop blocks
// while it is empty, and Close unblocks every waiter on both sides.
//
// The zero Queue is not usable; call NewQueue.
type Queue[T any] struct {
	// TODO: a mutex, one or more sync.Cond over it, a ring, a count, a closed flag.
}

// NewQueue returns an empty queue holding at most capacity items.
// It panics if capacity is not positive.
func NewQueue[T any](capacity int) *Queue[T] {
	panic("TODO: implement NewQueue")
}

// Push appends v, blocking while the queue is full. It returns ErrClosed if
// the queue is closed — whether it was already closed on entry, or was closed
// while this call was waiting for room.
func (q *Queue[T]) Push(v T) error {
	panic("TODO: implement Push")
}

// Pop removes and returns the oldest item, blocking while the queue is empty.
//
// A closed queue still yields what it already holds: Pop returns ErrClosed only
// once the queue is closed AND drained.
func (q *Queue[T]) Pop() (T, error) {
	panic("TODO: implement Pop")
}

// Close unblocks every waiter on both sides. It is safe to call more than once,
// and safe to call concurrently with Push and Pop.
func (q *Queue[T]) Close() {
	panic("TODO: implement Close")
}
