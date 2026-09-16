// Package abandoned holds five workers, each of which sends a result to a
// caller that has already given up, and a way to ask the runtime which of them
// are still there.
package abandoned

import (
	"runtime"
	"runtime/pprof"
	"strings"
	"time"
)

// Five workers, each abandoned the moment it is launched.
func UnbufferedWorker(ch chan int) { ch <- 1 }
func CapOneWorker(ch chan int)     { ch <- 1 }

func DefaultWorker(ch chan int) {
	select {
	case ch <- 1:
	default:
	}
}

func CapOneTwiceWorker(ch chan int) { ch <- 1; ch <- 2 }
func LateDrainWorker(ch chan int)   { ch <- 1 }

// Parked reports whether any goroutine is still inside a function whose name
// contains fn.
//
// A synctest bubble cannot be used for this: synctest.Test fails a test that
// ends with goroutines still running, and a goroutine still running is exactly
// what is being measured here. The goroutine profile names functions by full
// import path, hence the "."+fn match.
func Parked(fn string) bool {
	for range 3 {
		runtime.Gosched()
		time.Sleep(time.Millisecond)
	}
	var sb strings.Builder
	if err := pprof.Lookup("goroutine").WriteTo(&sb, 1); err != nil {
		return false
	}
	return strings.Contains(sb.String(), "."+fn)
}
