// Package drill — C10/11 model-swallows-errors.
package drill

import "errors"

// ErrFull is returned by Put when the store is at capacity and the key is new.
var ErrFull = errors.New("store full")

// Store is a map with a fixed number of keys. Updating a key it already holds
// always succeeds; adding a new key fails with ErrFull once it holds Cap keys.
type Store struct {
	capacity int
	m        map[int]int
}

// New returns an empty Store that holds at most capacity keys.
func New(capacity int) *Store {
	return &Store{capacity: capacity, m: make(map[int]int)}
}

// Put stores v under k.
func (s *Store) Put(k, v int) error {
	if _, ok := s.m[k]; !ok && len(s.m) >= s.capacity {
		return ErrFull
	}
	s.m[k] = v
	return nil
}

// Get returns the value stored under k.
func (s *Store) Get(k int) (int, bool) {
	v, ok := s.m[k]
	return v, ok
}

// Delete removes k. Deleting a key the store does not hold is a no-op.
func (s *Store) Delete(k int) { delete(s.m, k) }

// Len reports how many keys the store holds.
func (s *Store) Len() int { return len(s.m) }
