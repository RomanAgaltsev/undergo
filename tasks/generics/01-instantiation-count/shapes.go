// Package shapes instantiates one generic method and one generic function over
// the same five type pairs, so you can count what the compiler emitted.
//
// You do not implement anything here. Work out the answers, write them into
// prediction.yaml, then verify.
package shapes

// Box is a generic type with a generic method.
type Box[T any] struct{ V T }

// MapM is a generic METHOD: it declares its own type parameter R, on top of the
// T already bound by the receiver. Methods could not do this before Go 1.27.
func (b Box[T]) MapM[R any](f func(T) R) R { return f(b.V) }

// MapF is the equivalent generic FUNCTION: both type parameters are its own.
func MapF[T, R any](v T, f func(T) R) R { return f(v) }

// MyInt is a defined type whose underlying type is int. Whether that matters
// here is most of the task.
type MyInt int

// Callsites instantiates MapM and MapF over five type pairs each:
//
//	(int, string)   (int, int64)   (string, int)   (MyInt, string)   (*int, bool)
//
// It exists so that every instantiation is reachable and lands in the archive.
func Callsites() []any {
	return []any{
		Box[int]{1}.MapM(func(i int) string { return "" }),
		Box[int]{2}.MapM(func(i int) int64 { return 0 }),
		Box[string]{"a"}.MapM(func(s string) int { return 0 }),
		Box[MyInt]{3}.MapM(func(m MyInt) string { return "" }),
		Box[*int]{nil}.MapM(func(p *int) bool { return false }),

		MapF(1, func(i int) string { return "" }),
		MapF(2, func(i int) int64 { return 0 }),
		MapF("a", func(s string) int { return 0 }),
		MapF(MyInt(3), func(m MyInt) string { return "" }),
		MapF((*int)(nil), func(p *int) bool { return false }),
	}
}
