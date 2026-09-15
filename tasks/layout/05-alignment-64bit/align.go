// Package align asks what Go guarantees about where a field lands, and what it
// only happens to do.
//
// The three structs below hold the same two things — a bool and a 64-bit
// counter — declared three different ways.
package align

import (
	"structs"
	"sync/atomic"
	"unsafe"
)

// Plain puts a raw int64 after a bool. On a 32-bit platform this is the classic
// atomic-alignment bug waiting to happen.
type Plain struct {
	Flag bool
	N    int64
}

// Guarded uses atomic.Int64, which carries its own alignment requirement.
type Guarded struct {
	Flag bool
	N    atomic.Int64
}

// Hosted asks the compiler for host (C-compatible) layout rules.
type Hosted struct {
	_    structs.HostLayout
	Flag bool
	N    int64
}

// Int64Align reports the alignment the compiler gives a bare int64.
func Int64Align() uintptr { return unsafe.Alignof(int64(0)) }

// AtomicSize reports the size of an atomic.Int64.
func AtomicSize() uintptr { var g Guarded; return unsafe.Sizeof(g.N) }

// PlainOffset reports where N begins in Plain.
func PlainOffset() uintptr { var p Plain; return unsafe.Offsetof(p.N) }

// GuardedOffset reports where N begins in Guarded.
func GuardedOffset() uintptr { var g Guarded; return unsafe.Offsetof(g.N) }

// HostedSize reports the size of Hosted.
func HostedSize() uintptr { var h Hosted; return unsafe.Sizeof(h) }

// PlainSize reports the size of Plain, for comparison with Hosted.
func PlainSize() uintptr { var p Plain; return unsafe.Sizeof(p) }
