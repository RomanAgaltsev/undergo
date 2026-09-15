//go:build ignore

// VIOLATION: a uintptr round-trip through a struct field.
package patterns

import "unsafe"

type Holder struct{ Addr uintptr }

func V3(p *int) *Holder { return &Holder{Addr: uintptr(unsafe.Pointer(p))} }

func V3Back(h *Holder) *int { return (*int)(unsafe.Pointer(h.Addr)) }
