// Package gctrace parses the lines GODEBUG=gctrace=1 writes to stderr.
//
// The format is documented in the runtime package's own doc comment and is
// stable enough to parse, which is why every article about Go memory debugging
// starts here. It is also a debug format with no compatibility promise, which
// is the last question.
package gctrace

// Cycle is one collection.
type Cycle struct {
	Num        int     // the cycle number
	PercentCPU int     // cumulative percentage of CPU used by the collector
	HeapStart  int     // MB live at the start of the mark phase
	HeapEnd    int     // MB live at the end
	HeapGoal   int     // MB the collector was aiming for
	Procs      int     // GOMAXPROCS at the time
	Forced     bool    // triggered by runtime.GC rather than by the pacer
	WallMS     float64 // the three clock figures added together
}

// Parse reads gctrace output and returns one Cycle per collection.
//
// Lines that are not collections — scavenger output, build messages, anything
// else on the stream — are ignored rather than rejected. A malformed
// collection line is an error.
func Parse(trace string) ([]Cycle, error) {
	panic("implement Parse")
}

// TotalForced reports how many of the cycles were forced.
func TotalForced(cs []Cycle) int {
	panic("implement TotalForced")
}
