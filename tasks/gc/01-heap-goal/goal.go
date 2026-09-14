// Package goal reads the garbage collector's heap goal and compares it with the
// documented formula, at two heap sizes.
//
// You do not implement anything here. Work out the answers, write them into
// prediction.yaml, then verify.
package goal

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"runtime/metrics"
)

// Heap sizes, in KiB of retained 1 KiB allocations.
const (
	smallHeapKiB = 4 * 1024  // 4 MiB
	largeHeapKiB = 64 * 1024 // 64 MiB
)

// live is the retained set. It is package-level so it stays reachable for the
// whole measurement rather than depending on the caller's stack frame.
var live [][]byte

// read returns one runtime metric.
func read(name string) uint64 {
	s := []metrics.Sample{{Name: name}}
	metrics.Read(s)
	return s[0].Value.Uint64()
}

// ratio returns the heap goal divided by the live set, at the given GOGC and
// retained size.
func ratio(gogc, kib int) float64 {
	previous := debug.SetGCPercent(gogc)
	defer debug.SetGCPercent(previous)

	live = nil
	for range kib {
		live = append(live, make([]byte, 1024))
	}

	// Twice, and the second is the one that counts: the goal read after a
	// single collection is still based on the live set from before it.
	runtime.GC()
	runtime.GC()

	l := read("/gc/heap/live:bytes")
	g := read("/gc/heap/goal:bytes")
	return float64(g) / float64(l)
}

// Formula is the documented rule: a goal of live x (1 + GOGC/100), rounded to
// one decimal place.
func Formula(gogc int) string {
	return fmt.Sprintf("%.1f", 1+float64(gogc)/100)
}

// LargeHeapRatio returns the measured ratio over a 64 MiB live set, rounded to
// one decimal place so it is a value you can write down.
func LargeHeapRatio(gogc int) string {
	return fmt.Sprintf("%.1f", ratio(gogc, largeHeapKiB))
}

// SmallHeapExceedsFormula reports whether the measured ratio over a 4 MiB live
// set is above what the formula predicts, at every one of three GOGC settings.
func SmallHeapExceedsFormula() bool {
	for _, gogc := range []int{100, 200, 400} {
		predicted := 1 + float64(gogc)/100
		if ratio(gogc, smallHeapKiB) <= predicted {
			return false
		}
	}
	return true
}

// LargeHeapMatchesFormula reports whether the measured ratio over a 64 MiB live
// set is within 2% of the formula, at every one of three GOGC settings.
func LargeHeapMatchesFormula() bool {
	for _, gogc := range []int{100, 200, 400} {
		predicted := 1 + float64(gogc)/100
		got := ratio(gogc, largeHeapKiB)
		if got < predicted*0.98 || got > predicted*1.02 {
			return false
		}
	}
	return true
}

// GapShrinksAsHeapGrows reports whether the distance between the measured ratio
// and the formula is smaller on the large heap than on the small one.
func GapShrinksAsHeapGrows() bool {
	const gogc = 400
	predicted := 1 + float64(gogc)/100
	small := ratio(gogc, smallHeapKiB) - predicted
	large := ratio(gogc, largeHeapKiB) - predicted
	return large < small
}
