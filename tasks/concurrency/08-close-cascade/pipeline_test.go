package cascade

import (
	"context"
	"runtime"
	"runtime/pprof"
	"strings"
	"testing"
	"time"
)

// liveIn reports how many goroutines are currently inside a function of this
// package whose name contains fn. The goroutine profile names functions by full
// import path, hence the "."+fn match.
func liveIn(t *testing.T, fn string) int {
	t.Helper()
	var sb strings.Builder
	if err := pprof.Lookup("goroutine").WriteTo(&sb, 1); err != nil {
		t.Fatal(err)
	}
	return strings.Count(sb.String(), "."+fn)
}

// settle gives abandoned goroutines time to reach their blocking operation, so
// that a leak is visible rather than merely likely.
func drain(ch <-chan int) {
	for range ch { // deliberately discarding
		_ = 0
	}
}

// leakCheck captures a baseline now and returns the assertion to run later.
// It is used as `defer leakCheck(t, ...)()` so that nothing is assigned before
// the first call into the package under test — with a panicking stub, every
// statement after that call is unreachable, and a variable assigned before it
// and read after it reads to staticcheck as never used.
func leakCheck(t *testing.T, fns ...string) func() {
	t.Helper()
	before := make(map[string]int, len(fns))
	for _, fn := range fns {
		before[fn] = liveIn(t, fn)
	}
	return func() {
		settle()
		for _, fn := range fns {
			if got := liveIn(t, fn); got > before[fn] {
				t.Errorf("%s left %d goroutine(s) behind after cancellation", fn, got-before[fn])
			}
		}
	}
}

func settle() {
	for range 5 {
		runtime.Gosched()
		time.Sleep(2 * time.Millisecond)
	}
}

func TestLinearPipeline(t *testing.T) {
	ctx := context.Background()
	var got []int
	for v := range Square(ctx, Gen(ctx, 1, 2, 3, 4)) {
		got = append(got, v)
	}
	want := []int{1, 4, 9, 16}
	if len(got) != len(want) {
		t.Fatalf("got %d values; want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v; want %v", got, want)
		}
	}
}

func TestGenClosesItsChannel(t *testing.T) {
	ctx := context.Background()
	ch := Gen(ctx, 1)
	<-ch
	if _, ok := <-ch; ok {
		t.Fatal("Gen's channel yielded a second value; want it closed")
	}
}

// Every value survives the round trip through the graph, exactly once. Order
// across branches is not specified, so this compares multisets.
func TestFanOutFanInDeliversEveryValueOnce(t *testing.T) {
	const n = 40
	ctx := context.Background()

	vs := make([]int, n)
	for i := range vs {
		vs[i] = i
	}

	branches := FanOut(ctx, Gen(ctx, vs...), 4)
	squared := make([]<-chan int, len(branches))
	for i, b := range branches {
		squared[i] = Square(ctx, b)
	}

	counts := map[int]int{}
	for v := range FanIn(ctx, squared...) {
		counts[v]++
	}

	if len(counts) != n {
		t.Fatalf("received %d distinct values; want %d", len(counts), n)
	}
	for i := range n {
		if counts[i*i] != 1 {
			t.Fatalf("value %d received %d times; want 1", i*i, counts[i*i])
		}
	}
}

// FanIn's channel must close once every input is drained. If it does not, this
// range never ends and the test times out.
func TestFanInClosesWhenAllInputsClose(t *testing.T) {
	ctx := context.Background()
	a, b := Gen(ctx, 1, 2), Gen(ctx, 3, 4)
	seen := 0
	for range FanIn(ctx, a, b) {
		seen++
	}
	if seen != 4 {
		t.Fatalf("received %d values; want 4", seen)
	}
}

// The test that matters. Cancellation arrives mid-stream — after five of forty
// values — and every goroutine in the graph must unwind. A cancel test that
// cancels before starting, or after draining, proves nothing.
func TestCancelMidStreamLeaksNothing(t *testing.T) {
	const n = 40
	defer leakCheck(t, "Gen", "Square", "FanOut", "FanIn")()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	vs := make([]int, n)
	for i := range vs {
		vs[i] = i
	}

	branches := FanOut(ctx, Gen(ctx, vs...), 4)
	squared := make([]<-chan int, len(branches))
	for i, b := range branches {
		squared[i] = Square(ctx, b)
	}
	merged := FanIn(ctx, squared...)

	for range 5 {
		<-merged
	}
	cancel()

	// Drain whatever is already in flight; the channel must eventually close.
	drain(merged)
}

// Cancelling must close the downstream channel, not merely stop feeding it.
func TestCancelClosesDownstream(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	out := Square(ctx, Gen(ctx, 1, 2, 3, 4, 5, 6, 7, 8))
	<-out
	cancel()

	deadline := time.After(5 * time.Second)
	for {
		select {
		case _, ok := <-out:
			if !ok {
				return // closed, as required
			}
		case <-deadline:
			t.Fatal("Square's channel was never closed after cancellation")
		}
	}
}

func TestFanOutPanicsOnNonPositiveN(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("FanOut(ctx, in, 0) did not panic")
		}
	}()
	FanOut(context.Background(), Gen(context.Background()), 0)
}

// The consumer cancels and then simply walks away, without draining. This is
// the honest version of abandonment: nothing downstream will ever receive
// again, so every stage must unwind on ctx alone.
//
// Draining after cancellation — as TestCancelMidStreamLeaksNothing does — keeps
// the stages unblocked and therefore hides a stage whose send has no
// cancellation case. This test does not drain, and does not hide it.
func TestCancelWithoutDrainingLeaksNothing(t *testing.T) {
	const n = 40
	defer leakCheck(t, "Gen", "Square", "FanOut", "FanIn")()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	vs := make([]int, n)
	for i := range vs {
		vs[i] = i
	}

	branches := FanOut(ctx, Gen(ctx, vs...), 4)
	squared := make([]<-chan int, len(branches))
	for i, b := range branches {
		squared[i] = Square(ctx, b)
	}
	merged := FanIn(ctx, squared...)

	for range 5 {
		<-merged
	}
	cancel()
	// No drain. merged is abandoned exactly where it is.
}

// The minimal shape that pins a stage's send-side cancellation.
//
// One stage, one consumer, no fan-out to unwind through. The consumer takes a
// single value, cancels, and never receives again. Square is now holding a value
// nobody will ever take, and the only thing that can free it is ctx. A Square
// whose send is a bare `out <- v*v` rather than a select stays parked forever.
//
// The wider graph cannot show this: there, cancellation unwinds from the
// producer down, and each stage's input closes in turn, so a stage reaches the
// end of its range loop before its send ever becomes permanent.
func TestStageSendUnblocksOnCancelAlone(t *testing.T) {
	defer leakCheck(t, "Square")()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	vs := make([]int, 40)
	for i := range vs {
		vs[i] = i
	}
	out := Square(ctx, Gen(ctx, vs...))

	<-out // take exactly one, leaving Square mid-stream
	cancel()
	// No drain, ever.
}
