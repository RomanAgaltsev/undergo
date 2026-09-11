// Package drill — C11/06 generic-set.
package drill

// Set is a collection of unique elements that also remembers insertion order.
type Set[T comparable] struct {
	seen  map[T]struct{}
	order []T
}

// NewSet builds an empty set.
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{seen: make(map[T]struct{})}
}

// Add inserts v if not already present.
func (s *Set[T]) Add(v T) {
	if _, ok := s.seen[v]; ok {
		return
	}
	s.seen[v] = struct{}{}
	s.order = append(s.order, v)
}

// Contains reports whether v is in the set.
func (s *Set[T]) Contains(v T) bool {
	for _, e := range s.order {
		if e == v {
			return true
		}
	}
	return false
}

// Items returns the elements in insertion order.
func (s *Set[T]) Items() []any {
	out := make([]any, 0, len(s.order))
	for _, e := range s.order {
		out = append(out, e)
	}
	return out
}

// First returns the first-inserted element.
func (s *Set[T]) First() T {
	var zero T
	if len(s.order) == 0 {
		return zero
	}
	return s.order[0]
}
