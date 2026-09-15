//go:build ignore

// Pattern 3: arithmetic, folded into the conversion.
package patterns

import "unsafe"

type S struct{ A, B int64 }

func P2(s *S) *int64 {
	return (*int64)(unsafe.Pointer(uintptr(unsafe.Pointer(s)) + unsafe.Offsetof(s.B)))
}
