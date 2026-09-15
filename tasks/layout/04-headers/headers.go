// Package headers asks how large Go's built-in types really are, and where
// their data actually lives.
//
// Every answer here is a compile-time constant. Nothing is measured at run
// time, and nothing depends on how much data the values hold.
package headers

import "unsafe"

// Header nests one of each of the types people guess the size of.
type Header struct {
	S   string
	Sl  []int
	M   map[int]int
	Ch  chan int
	F   func()
	I   any
	Ptr *int
}

// Sizes reports unsafe.Sizeof for each field and for the whole struct.
func Sizes() map[string]uintptr {
	var h Header
	return map[string]uintptr{
		"string": unsafe.Sizeof(h.S),
		"slice":  unsafe.Sizeof(h.Sl),
		"map":    unsafe.Sizeof(h.M),
		"chan":   unsafe.Sizeof(h.Ch),
		"func":   unsafe.Sizeof(h.F),
		"iface":  unsafe.Sizeof(h.I),
		"ptr":    unsafe.Sizeof(h.Ptr),
		"total":  unsafe.Sizeof(h),
	}
}

// OffsetOfMap reports where the map field begins.
func OffsetOfMap() uintptr {
	var h Header
	return unsafe.Offsetof(h.M)
}

// BigSliceSize reports Sizeof for a slice holding a million elements, for
// comparison with the empty one.
func BigSliceSize() uintptr {
	s := make([]int, 1_000_000)
	return unsafe.Sizeof(s)
}
