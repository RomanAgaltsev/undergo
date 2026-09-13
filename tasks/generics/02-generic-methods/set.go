// Package set holds a typed collection whose Map method changes the element
// type — which a method could not do before Go 1.27.
//
// Before generic methods, a transformation that introduces a new type parameter
// had to be a package-level function taking the collection as an argument,
// because only functions could declare type parameters. Here it is a method.
package set

import (
	"context"
	"fmt"
)

// Sink keeps results reachable so the allocation test measures the work rather
// than dead code elimination.
//
// It is []int rather than any on purpose: assigning a slice to an interface
// boxes the slice header, which is an allocation of the test's own making.
var Sink []int

// MapError reports which item failed, and wraps the error it failed with.
type MapError struct {
	Index int
	Err   error
}

func (e *MapError) Error() string { return fmt.Sprintf("item %d: %v", e.Index, e.Err) }

// Unwrap lets errors.Is and errors.As reach the underlying error.
func (e *MapError) Unwrap() error { return e.Err }

// Set is an ordered collection of T.
type Set[T any] struct {
	items []T
}

// New returns a Set holding items, in order.
func New[T any](items ...T) *Set[T] {
	return &Set[T]{items: items}
}

// Len returns how many items the Set holds.
func (s *Set[T]) Len() int { return len(s.items) }

// Map applies fn to every item and returns the results in input order.
//
// R is the method's own type parameter: the receiver already bound T, and Map
// introduces a second parameter that has nothing to do with it. Callers write no
// type arguments — R is inferred from fn.
//
// The contract the frozen tests enforce:
//
//   - results come back in input order, one per item;
//   - an empty Set returns a nil slice and a nil error, and never calls fn;
//   - ctx is checked before each item, so a context already cancelled on entry
//     returns ctx.Err() without calling fn at all;
//   - the first error stops the walk — later items are not passed to fn — and is
//     returned as a *MapError carrying the index of the item that failed;
//   - exactly one slice is allocated, sized up front.
func (s *Set[T]) Map[R any](ctx context.Context, fn func(context.Context, T) (R, error)) ([]R, error) {
	panic("undergo: implement Map")
}
