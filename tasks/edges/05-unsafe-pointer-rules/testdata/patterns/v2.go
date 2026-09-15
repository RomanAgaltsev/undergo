//go:build ignore

// VIOLATION: pointer arithmetic that walks past the end of the allocation.
package patterns

import "unsafe"

func V2(s []int32) *int32 {
	base := unsafe.Pointer(&s[0])
	return (*int32)(unsafe.Add(base, uintptr(len(s)+16)*unsafe.Sizeof(s[0])))
}
