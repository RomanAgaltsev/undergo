package tree

import (
	"runtime"
	"slices"
	"testing"
)

func build(vs ...int) *Tree {
	var t *Tree
	for _, v := range vs {
		t = Insert(t, v)
	}
	return t
}

var sample = []int{50, 30, 70, 20, 40, 60, 80, 10}

func sorted() []int {
	s := slices.Clone(sample)
	slices.Sort(s)
	return s
}

// TestAllIsInOrder checks the push sequence visits every value in order.
func TestAllIsInOrder(t *testing.T) {
	var got []int
	for v := range All(build(sample...)) {
		got = append(got, v)
	}
	if !slices.Equal(got, sorted()) {
		t.Errorf("All did not produce the values in order")
	}
}

// TestAllStopsDescending checks the recursion propagates yield's result. A walk
// that ignores it still visits the whole tree, and this is the test that
// catches it.
func TestAllStopsDescending(t *testing.T) {
	visited := 0
	for range All(build(sample...)) {
		visited++
		break
	}
	if visited != 1 {
		t.Fatalf("the loop body ran %d times after one break", visited)
	}
}

// TestAllEmptyTree checks the nil tree yields nothing rather than panicking.
func TestAllEmptyTree(t *testing.T) {
	for range All(nil) {
		t.Fatal("an empty tree yielded a value")
	}
}

// TestWalkPullsInOrder checks the pull sequence returns the same values.
func TestWalkPullsInOrder(t *testing.T) {
	next, stop := Walk(build(sample...))
	defer stop()

	var got []int
	for {
		v, ok := next()
		if !ok {
			break
		}
		got = append(got, v)
	}
	if !slices.Equal(got, sorted()) {
		t.Errorf("Walk did not produce the values in order")
	}
}

// TestWalkResumes is the point of a pull sequence: take a few, do something
// else, then carry on from where you were.
func TestWalkResumes(t *testing.T) {
	next, stop := Walk(build(sample...))
	defer stop()

	want := sorted()
	for i := range 3 {
		if v, ok := next(); !ok || v != want[i] {
			t.Fatalf("value %d of the first run was wrong", i)
		}
	}
	for i := 3; i < len(want); i++ {
		if v, ok := next(); !ok || v != want[i] {
			t.Fatalf("value %d after resuming was wrong", i)
		}
	}
}

// TestStopIsIdempotentAndFinal checks the two properties the doc comment asks
// for: stop may be called twice, and next reports false afterwards.
func TestStopIsIdempotentAndFinal(t *testing.T) {
	next, stop := Walk(build(sample...))
	if _, ok := next(); !ok {
		t.Fatal("the first value was missing")
	}

	stop()
	stop() // must not panic

	if _, ok := next(); ok {
		t.Error("next returned a value after stop")
	}
}

// TestStopReleasesTheWalk checks that stopping actually releases what the pull
// sequence is holding. Twenty walks are started, half-consumed and stopped; if
// stop does nothing, the goroutines they are running stay alive.
func TestStopReleasesTheWalk(t *testing.T) {
	before := runtime.NumGoroutine()

	for range 20 {
		next, stop := Walk(build(sample...))
		next()
		next()
		stop()
	}

	runtime.GC()
	runtime.Gosched()

	if after := runtime.NumGoroutine(); after > before+2 {
		t.Errorf("stopping twenty half-consumed walks left goroutines behind")
	}
}
