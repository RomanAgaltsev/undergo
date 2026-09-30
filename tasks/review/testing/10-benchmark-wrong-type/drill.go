// Package drill — C10/10 benchmark-wrong-type.
package drill

// Queue is a FIFO of ints.
type Queue interface {
	Push(v int)
	Pop() (int, bool)
}

// SliceQueue is the queue in use today: Push appends and Pop reslices.
type SliceQueue struct {
	items []int
}

// NewSliceQueue returns an empty SliceQueue.
func NewSliceQueue() *SliceQueue { return &SliceQueue{} }

// Push adds v at the back.
func (q *SliceQueue) Push(v int) { q.items = append(q.items, v) }

// Pop removes and returns the front value, or reports false if q is empty.
func (q *SliceQueue) Pop() (int, bool) {
	if len(q.items) == 0 {
		return 0, false
	}
	v := q.items[0]
	q.items = q.items[1:]
	return v, true
}

// RingQueue replaces SliceQueue: a ring buffer that stops allocating once it
// has grown to the queue's working size.
type RingQueue struct {
	buf        []int
	head, size int
}

// NewRingQueue returns an empty RingQueue with room for capacity values.
func NewRingQueue(capacity int) *RingQueue {
	return &RingQueue{buf: make([]int, max(capacity, 1))}
}

// Push adds v at the back, doubling the ring if it is full.
func (q *RingQueue) Push(v int) {
	if q.size == len(q.buf) {
		q.grow()
	}
	q.buf[(q.head+q.size)%len(q.buf)] = v
	q.size++
}

// Pop removes and returns the front value, or reports false if q is empty.
func (q *RingQueue) Pop() (int, bool) {
	if q.size == 0 {
		return 0, false
	}
	v := q.buf[q.head]
	q.head = (q.head + 1) % len(q.buf)
	q.size--
	return v, true
}

func (q *RingQueue) grow() {
	next := make([]int, 2*len(q.buf))
	for i := range q.size {
		next[i] = q.buf[(q.head+i)%len(q.buf)]
	}
	q.buf, q.head = next, 0
}
