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
	"slices"
	"strings"
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

// Ordering ranks the four settings by how many collections each performed, most
// first, joined by ">".
func Ordering(garbageMiB int) string {
	type result struct {
		name   string
		cycles uint64
	}
	results := []result{
		{"gogc50", Cycles(50, garbageMiB)},
		{"gogc100", Cycles(100, garbageMiB)},
		{"gogc400", Cycles(400, garbageMiB)},
		{"gogc800", Cycles(800, garbageMiB)},
	}
	slices.SortStableFunc(results, func(a, b result) int {
		switch {
		case a.cycles > b.cycles:
			return -1
		case a.cycles < b.cycles:
			return 1
		default:
			return 0
		}
	})

	names := make([]string, 0, len(results))
	for _, r := range results {
		names = append(names, r.name)
	}
	return strings.Join(names, ">")
}
