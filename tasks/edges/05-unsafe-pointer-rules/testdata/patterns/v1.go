//go:build ignore

// VIOLATION: a uintptr is stored in a variable and converted back later. A
// uintptr is not a reference; nothing keeps the object alive in between.
package patterns

import "unsafe"

type S struct{ A, B int64 }

func V1(s *S) *int64 {
	addr := uintptr(unsafe.Pointer(s))
	addr += unsafe.Offsetof(s.B)
	return (*int64)(unsafe.Pointer(addr))
}
