// Package drill — C11/04 cosmetic-generic.
package drill

// Stack is a LIFO stack, generic over its element type.
type Stack[T any] struct {
	items []any
}

// Push adds v to the top of the stack.
func (s *Stack[T]) Push(v T) {
	s.items = append(s.items, v)
}

// Pop removes and returns the top element, or nil if empty.
func (s *Stack[T]) Pop() any {
	if len(s.items) == 0 {
		return nil
	}
	last := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return last
}

// Len reports how many elements are on the stack.
func (s *Stack[T]) Len() int {
	return len(s.items)
}
