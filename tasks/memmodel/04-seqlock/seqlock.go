// Package seqlock is a sequence lock: one writer, many readers, and readers
// that never block the writer and never take a lock.
//
// A reader reads the sequence number, reads the payload, reads the sequence
// number again, and retries if it changed or was odd. Implement it so that a
// reader never returns a value assembled from two different writes, and so that
// the race detector stays silent.
package seqlock

// Point is the payload. It is deliberately wider than one word, so a torn read
// is possible and detectable.
type Point struct{ X, Y, Z int64 }

// Seqlock guards a Point.
type Seqlock struct {
	// Add the fields you need. A sequence number and the payload is enough,
	// but their types are the exercise.
}

// Write publishes p. Only one goroutine calls Write.
func (s *Seqlock) Write(p Point) {
	panic("implement Write")
}

// Read returns the most recently published Point, retrying while a write is in
// progress. Many goroutines call Read concurrently.
func (s *Seqlock) Read() Point {
	panic("implement Read")
}
