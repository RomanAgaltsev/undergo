// Package bce holds six functions that index a slice in six different shapes.
//
// You do not implement anything here. Work out which keep a bounds check, write
// your answers into prediction.yaml, then verify.
package bce

// Ranged sums with a range loop.
func Ranged(xs []int) int {
	t := 0
	for _, x := range xs {
		t += x
	}
	return t
}

// Counted sums with an index against len.
func Counted(xs []int) int {
	t := 0
	for i := 0; i < len(xs); i++ {
		t += xs[i]
	}
	return t
}

// TwoSlices walks one slice and indexes another with the same counter.
func TwoSlices(xs, ys []int) int {
	t := 0
	for i := range xs {
		t += ys[i]
	}
	return t
}

// Unchecked indexes with a value it was handed.
func Unchecked(xs []int, i int) int { return xs[i] }

// AfterOneCheck reads four elements after a single length test.
func AfterOneCheck(xs []int) int {
	if len(xs) < 4 {
		return 0
	}
	return xs[0] + xs[1] + xs[2] + xs[3]
}

// NextElement reads the element after the current one.
func NextElement(xs []int) int {
	t := 0
	for i := 0; i < len(xs)-1; i++ {
		t += xs[i+1]
	}
	return t
}
