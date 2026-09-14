// Package cycles runs the same allocation workload at four GOGC settings and
// counts the collections each one performed.
//
// You do not implement anything here. Work out the answers, write them into
// prediction.yaml, then verify.
//
// Read the README before the code: this task does not ask you for a number of
// collections, and the reason it does not is the thing worth learning.
package cycles

import (
	"runtime"
	"runtime/debug"
	"runtime/metrics"
)

// live is retained for the whole run, so the workload has a stable live set.
var live [][]byte

// Garbage is overwritten every iteration, so each allocation dies immediately.
//
// It is exported so that the compiler cannot decide the allocations are dead
// and remove the workload this task exists to measure.
var Garbage []byte

func readCycles() uint64 {
	s := []metrics.Sample{{Name: "/gc/cycles/total:gc-cycles"}}
	metrics.Read(s)
	return s[0].Value.Uint64()
}

// Cycles counts the collections performed while allocating garbageMiB of
// immediately-dead objects at the given GOGC, over a 2 MiB live set.
func Cycles(gogc int, garbageMiB int) uint64 {
	previous := debug.SetGCPercent(gogc)
	defer debug.SetGCPercent(previous)

	live = nil
	for range 2048 {
		live = append(live, make([]byte, 1024))
	}
	// The first collection establishes the live set; the second settles the
	// goal on it. Counting starts after both.
	runtime.GC()
	runtime.GC()

	before := readCycles()
	for range garbageMiB * 1024 {
		Garbage = make([]byte, 1024)
	}
	return readCycles() - before
}

// Beats reports whether the first setting performed more collections than the
// second.
//
// Pairwise rather than a full ranking: the counts for adjacent settings can sit
// close enough that load on the machine reorders them, while settings an order
// of magnitude apart do not move. See EXPLANATION.md.
func Beats(gogcA, gogcB, garbageMiB int) bool {
	return Cycles(gogcA, garbageMiB) > Cycles(gogcB, garbageMiB)
}
