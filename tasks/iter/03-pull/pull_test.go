package pull

import (
	"iter"
	"runtime"
	"slices"
	"testing"
	"time"
)

// tree builds a small BST holding 1..7.
func tree() *Node {
	return &Node{
		Value: 4,
		Left:  &Node{Value: 2, Left: &Node{Value: 1}, Right: &Node{Value: 3}},
		Right: &Node{Value: 6, Left: &Node{Value: 5}, Right: &Node{Value: 7}},
	}
}

// tracked returns a sequence over values and a flag reporting whether the
// sequence function ever returned — which it does when it is stopped.
type tracked struct {
	finished bool
}

func (tr *tracked) seq(values ...int) iter.Seq[int] {
	return func(yield func(int) bool) {
		defer func() { tr.finished = true }()
		for _, v := range values {
			if !yield(v) {
				return
			}
		}
	}
}

func TestInOrderWalksSorted(t *testing.T) {
	got := slices.Collect(tree().InOrder())
	if want := []int{1, 2, 3, 4, 5, 6, 7}; !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestInOrderOnAnEmptyTree(t *testing.T) {
	var empty *Node
	if got := slices.Collect(empty.InOrder()); len(got) != 0 {
		t.Errorf("got %v, want nothing", got)
	}
}

// TestInOrderStopsEarly catches the recursive walk that ignores yield's result:
// it collects the right values and keeps walking after the consumer stopped.
func TestInOrderStopsEarly(t *testing.T) {
	var seen []int
	for v := range tree().InOrder() {
		seen = append(seen, v)
		if v == 3 {
			break
		}
	}
	if want := []int{1, 2, 3}; !slices.Equal(seen, want) {
		t.Errorf("saw %v, want %v", seen, want)
	}
}

func TestMergeIsSorted(t *testing.T) {
	a := slices.Values([]int{1, 4, 6, 9})
	b := slices.Values([]int{2, 3, 7})

	got := slices.Collect(Merge(a, b))
	if want := []int{1, 2, 3, 4, 6, 7, 9}; !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestMergeHandlesEmptySides(t *testing.T) {
	full := []int{1, 2, 3}

	got := slices.Collect(Merge(slices.Values(full), slices.Values([]int(nil))))
	if !slices.Equal(got, full) {
		t.Errorf("empty on the right: got %v, want %v", got, full)
	}

	got = slices.Collect(Merge(slices.Values([]int(nil)), slices.Values(full)))
	if !slices.Equal(got, full) {
		t.Errorf("empty on the left: got %v, want %v", got, full)
	}

	if got := slices.Collect(Merge(slices.Values([]int(nil)), slices.Values([]int(nil)))); len(got) != 0 {
		t.Errorf("both empty: got %v, want nothing", got)
	}
}

func TestMergeKeepsDuplicates(t *testing.T) {
	got := slices.Collect(Merge(slices.Values([]int{1, 2, 2}), slices.Values([]int{2, 3})))
	if want := []int{1, 2, 2, 2, 3}; !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

// TestMergeStopsBothSources is the test a missing stop() fails: each source
// records whether its own function returned.
func TestMergeStopsBothSources(t *testing.T) {
	var ta, tb tracked

	taken := 0
	for range Merge(ta.seq(1, 3, 5, 7, 9), tb.seq(2, 4, 6, 8)) {
		taken++
		if taken == 2 {
			break
		}
	}

	// Stopping a pulled sequence unwinds it on its own goroutine, so give it
	// a moment to finish before looking.
	settle()

	if !ta.finished {
		t.Error("the first source was never stopped")
	}
	if !tb.finished {
		t.Error("the second source was never stopped")
	}
}

// TestMergeDoesNotLeakGoroutines is the same defect measured from the outside.
//
// It counts goroutines rather than detecting races — this repository has no race
// gate, and a leaked iter.Pull goroutine is not a race anyway. It is a goroutine
// parked forever, waiting to hand over an element nobody will ask for.
func TestMergeDoesNotLeakGoroutines(t *testing.T) {
	settle()
	before := runtime.NumGoroutine()

	for range 20 {
		var ta, tb tracked
		taken := 0
		for range Merge(ta.seq(1, 3, 5, 7, 9), tb.seq(2, 4, 6, 8)) {
			taken++
			if taken == 2 {
				break
			}
		}
	}

	settle()
	if after := runtime.NumGoroutine(); after > before {
		t.Errorf("goroutine count went from %d to %d — something was not stopped", before, after)
	}
}

// settle gives stopped pull goroutines a chance to exit before they are counted.
func settle() {
	for range 20 {
		runtime.Gosched()
		time.Sleep(time.Millisecond)
	}
}
