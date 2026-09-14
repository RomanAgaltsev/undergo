// Package zero is about allocations that occupy no space.
//
// You do not implement anything here. Work out the answers, write them into
// prediction.yaml, then verify.
package zero

import "unsafe"

// Sinks keep the allocations reachable so the compiler cannot fold them away.
var (
	P1 *struct{}
	P2 *struct{}
	A1 *[0]int
	A2 *[0]int
)

// Trailing has a zero-size field at the end, which is the case the spec singles
// out.
type Trailing struct {
	N int
	Z struct{}
}

// Leading has the zero-size field first.
type Leading struct {
	Z struct{}
	N int
}

// TwoEmptyStructsShareAnAddress reports whether two separate allocations of a
// zero-size struct compare equal.
func TwoEmptyStructsShareAnAddress() bool {
	P1 = new(struct{})
	P2 = new(struct{})
	return P1 == P2
}

// TwoZeroArraysShareAnAddress asks the same of a zero-length array.
func TwoZeroArraysShareAnAddress() bool {
	A1 = new([0]int)
	A2 = new([0]int)
	return A1 == A2
}

// TrailingFieldAddsPadding reports whether a trailing zero-size field makes the
// struct larger than the field before it.
func TrailingFieldAddsPadding() bool {
	return unsafe.Sizeof(Trailing{}) > unsafe.Sizeof(int(0))
}

// LeadingFieldAddsPadding asks the same of a leading zero-size field.
func LeadingFieldAddsPadding() bool {
	return unsafe.Sizeof(Leading{}) > unsafe.Sizeof(int(0))
}

// ZeroLengthSlicesOfOneArrayShareAData reports whether two zero-length slices
// taken from the same array have the same data pointer.
func ZeroLengthSlicesOfOneArrayShareAData() bool {
	var a [8]int
	s1 := a[0:0]
	s2 := a[4:4]
	return unsafe.SliceData(s1) == unsafe.SliceData(s2)
}

// EmptySliceAndNilSliceAreEqual reports whether an empty non-nil slice compares
// equal to nil.
func EmptySliceAndNilSliceAreEqual() bool {
	s := []int{}
	return s == nil
}
