// Package fold folds a slice with the element type's own Add method.
//
// Every line of this file was illegal at some point in Go's history, and two of
// them were illegal as recently as Go 1.25. If you are sure one of them does not
// compile, check which release you are remembering.
package fold

import "errors"

// Sink keeps results reachable so the allocation test measures the work rather
// than dead code elimination.
var Sink *Money

// ErrEmpty is returned when there is nothing to fold.
//
// There is no zero value to start from: the constraint promises Add and nothing
// else, so the function cannot invent an identity element.
var ErrEmpty = errors.New("fold: no items")

// Adder refers to itself in its own type parameter list.
//
// Before Go 1.26 this was rejected — "a generic type may not refer to itself in
// its type parameter list" — and the workaround cost an extra type parameter.
// It describes the types closed under their own addition: whatever A is, Add
// takes an A and gives an A back.
type Adder[A Adder[A]] interface {
	Add(A) A
}

// Money is the type the tests fold. Note that Add is on the value receiver, so
// Money itself satisfies Adder[Money].
type Money struct {
	Cents int64
}

// Add returns the sum of two amounts.
func (m Money) Add(other Money) Money { return Money{Cents: m.Cents + other.Cents} }

// SumAll folds items with Add and returns a pointer to the result.
//
// The contract the frozen tests enforce:
//
//   - the fold runs left to right, starting from the first element;
//   - an empty slice returns a nil pointer and ErrEmpty;
//   - a one-element slice returns a pointer to that element;
//   - exactly one allocation happens, and it is the result.
//
// Allocate the result with new taking an EXPRESSION, which is the other thing
// Go 1.26 added. No test can tell the two spellings apart — this one is on you,
// and written question 2 asks you to explain the difference.
func SumAll[A Adder[A]](items []A) (*A, error) {
	panic("undergo: implement SumAll")
}
