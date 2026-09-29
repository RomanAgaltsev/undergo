// Package varmake asks when a make whose length is not a constant allocates.
//
// Until Go 1.24 the answer was "every time": the compiler could not size a
// stack slot for a length it did not know, so make([]T, n) with a variable n
// went to the heap however small n turned out to be. Go 1.25 changed that.
// Both functions below make a slice that never leaves their frame; the only
// difference between them is the element type.
package varmake

// Bytes makes a byte slice of length n and uses it locally.
//
//go:noinline
func Bytes(n int) int {
	b := make([]byte, n)
	if n > 0 {
		b[0] = 1
	}
	return len(b)
}

// Words makes an int64 slice of length n and uses it locally.
//
//go:noinline
func Words(n int) int {
	w := make([]int64, n)
	if n > 0 {
		w[0] = 1
	}
	return len(w)
}
