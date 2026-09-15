// Package compare asks which interface comparisons the runtime refuses.
//
// Comparing two interface values compares their dynamic types first, and then,
// only if those match, their dynamic values. The second step is where it can go
// wrong: some types are not comparable at all, and the compiler cannot see
// which type an interface holds.
package compare

// WithSlice is uncomparable because one of its fields is.
type WithSlice struct{ S []int }

// Plain is comparable.
type Plain struct{ A, B int }

// Panics reports whether comparing a and b panicked.
func Panics(a, b any) (panicked bool) {
	defer func() {
		if recover() != nil {
			panicked = true
		}
	}()
	_ = a == b
	return false
}

// Equal reports whether a and b compare equal, treating a panic as false.
func Equal(a, b any) bool {
	if Panics(a, b) {
		return false
	}
	return a == b
}
