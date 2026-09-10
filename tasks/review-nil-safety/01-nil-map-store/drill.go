// Package drill — C2/01 nil-map-store.
package drill

// Store keeps integer counters by name.
type Store struct {
	data map[string]int
}

// NewStore returns a ready Store.
func NewStore() *Store {
	return &Store{}
}

// Set records v for key k.
func (s *Store) Set(k string, v int) {
	s.data[k] = v
}

// Get returns the value for k, or 0 if absent.
func (s *Store) Get(k string) int {
	return s.data[k]
}
