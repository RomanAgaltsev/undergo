// Package limit asks which of GOGC and GOMEMLIMIT is deciding the heap goal.
//
// You do not implement anything here. Work out the answers, write them into
// prediction.yaml, then verify.
//
// Every function restores the previous GOGC and memory limit before returning.
// Both are process-global, and a task that left a 16 MiB limit set would break
// every test that ran after it in the same binary.
package limit

import (
	"math"
	"runtime"
	"runtime/debug"
	"runtime/metrics"
)

// Live is the retained set, exported so the compiler cannot decide it is dead.
var Live [][]byte

// noLimit is what debug.SetMemoryLimit uses to mean "no limit".
const noLimit = math.MaxInt64

func read(name string) uint64 {
	s := []metrics.Sample{{Name: name}}
	metrics.Read(s)
	return s[0].Value.Uint64()
}

// retain fills Live with about mib MiB and settles the goal over two cycles.
func retain(mib int) {
	Live = nil
	for range mib * 1024 {
		Live = append(Live, make([]byte, 1024))
	}
	runtime.GC()
	runtime.GC()
}

// ClampedByLimit reports whether the memory limit, rather than GOGC, is what
// decided the heap goal.
//
// GOGC alone would ask for live x (1 + GOGC/100). If the real goal is well below
// that, something else is holding it down, and the only candidate here is the
// limit. The 10% margin exists so that ordinary overhead never decides the
// answer — see the explanation for why the configurations are chosen to sit
// nowhere near it.
func ClampedByLimit(gogc int, limitMiB int64) bool {
	previousGC := debug.SetGCPercent(gogc)
	defer debug.SetGCPercent(previousGC)
	previousLimit := debug.SetMemoryLimit(limitMiB << 20)
	defer debug.SetMemoryLimit(previousLimit)

	retain(8)

	live := read("/gc/heap/live:bytes")
	goal := read("/gc/heap/goal:bytes")

	if gogc < 0 {
		// GOGC is off, so there is no proportional goal at all: anything
		// finite must have come from the limit.
		return goal < uint64(noLimit)
	}
	unclamped := live * uint64(100+gogc) / 100
	return goal < unclamped*9/10
}

// LimitCanBeExceeded reports whether a program can hold more than the memory
// limit allows.
func LimitCanBeExceeded() bool {
	previousGC := debug.SetGCPercent(-1) // off, so only the limit can collect
	defer debug.SetGCPercent(previousGC)
	previousLimit := debug.SetMemoryLimit(32 << 20)
	defer debug.SetMemoryLimit(previousLimit)

	// Retain comfortably more than the limit permits.
	retain(64)

	total := read("/memory/classes/total:bytes")
	return total > 32<<20
}
