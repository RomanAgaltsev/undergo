package abandoned

import (
	"runtime"
	"runtime/pprof"
	"strings"
	"testing"
	"time"

	"github.com/RomanAgaltsev/undergo/internal/predict"
)

// parked reports whether any goroutine in the profile is still inside fn.
//
// This task cannot use a synctest bubble, because synctest.Test fails a test
// that ends with goroutines still running — and a goroutine still running is
// precisely this task's subject. The goroutine profile names functions by full
// import path, so the match is on "."+fn.
func parked(t *testing.T, fn string) bool {
	t.Helper()
	for range 3 {
		runtime.Gosched()
		time.Sleep(time.Millisecond)
	}
	var sb strings.Builder
	if err := pprof.Lookup("goroutine").WriteTo(&sb, 1); err != nil {
		t.Fatal(err)
	}
	return strings.Contains(sb.String(), "."+fn)
}

// Five workers, each abandoned by a caller that has already given up.
func unbufferedWorker(ch chan int) { ch <- 1 }
func capOneWorker(ch chan int)     { ch <- 1 }
func defaultWorker(ch chan int) {
	select {
	case ch <- 1:
	default:
	}
}
func capOneTwiceWorker(ch chan int) { ch <- 1; ch <- 2 }
func lateDrainWorker(ch chan int)   { ch <- 1 }

func TestPredictions(t *testing.T) {
	got := map[string]any{}

	ch1 := make(chan int)
	go unbufferedWorker(ch1)
	got["unbuffered_leaks"] = parked(t, "unbufferedWorker")

	ch2 := make(chan int, 1)
	go capOneWorker(ch2)
	got["capacity_one_leaks"] = parked(t, "capOneWorker")

	ch3 := make(chan int)
	go defaultWorker(ch3)
	got["select_with_default_leaks"] = parked(t, "defaultWorker")

	ch4 := make(chan int, 1)
	go capOneTwiceWorker(ch4)
	got["capacity_one_two_results_leaks"] = parked(t, "capOneTwiceWorker")

	ch5 := make(chan int)
	go lateDrainWorker(ch5)
	go func() { <-ch5 }()
	got["late_drain_leaks"] = parked(t, "lateDrainWorker")

	predict.Check(t, got)
}
