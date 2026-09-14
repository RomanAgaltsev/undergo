package intern

import (
	"fmt"
	"runtime"
	"testing"
)

// Sink keeps a measurement's data reachable across a heap reading.
var Sink any

func TestEqualStringsGiveEqualHandles(t *testing.T) {
	// Built two different ways so they cannot share a backing array: a
	// handle comparison must not degenerate into a pointer comparison of
	// the inputs.
	a := Intern("hello world")
	b := Intern("hello " + string([]byte("world")))

	if a != b {
		t.Error("handles for equal strings do not compare equal")
	}
}

func TestDistinctStringsGiveDistinctHandles(t *testing.T) {
	if Intern("alpha") == Intern("beta") {
		t.Error("handles for different strings compare equal")
	}
}

func TestValueRoundTrips(t *testing.T) {
	for _, s := range []string{"", "a", "hello world"} {
		if got := Value(Intern(s)); got != s {
			t.Errorf("Value(Intern(%q)) = %q", s, got)
		}
	}
}

func TestInternAllPreservesOrder(t *testing.T) {
	in := []string{"a", "b", "a", "c"}
	got := InternAll(in)

	if len(got) != len(in) {
		t.Fatalf("got %d handles, want %d", len(got), len(in))
	}
	for i := range in {
		if Value(got[i]) != in[i] {
			t.Errorf("handle %d holds %q, want %q", i, Value(got[i]), in[i])
		}
	}
	if got[0] != got[2] {
		t.Error("the two \"a\" entries did not intern to the same handle")
	}
}

// heapAfter returns the live heap once build's result is reachable.
func heapAfter(build func() any) uint64 {
	Sink = nil
	runtime.GC()

	Sink = build()

	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.HeapAlloc
}

// TestInterningRetainsLessThanStoringStrings measures both approaches in the
// same run and compares them against each other.
//
// It deliberately does not assert a byte count: retained bytes depend on the
// allocator, the platform and what else the test binary has done, and a hard
// number here would be the one flaky assertion in this track. The ratio is the
// durable claim.
func TestInterningRetainsLessThanStoringStrings(t *testing.T) {
	const (
		records = 200_000
		labels  = 16
	)

	// Each record's label is one of `labels` values, built fresh so that no
	// two records share a backing array.
	makeLabel := func(i int) string {
		return fmt.Sprintf("service-name-label-%02d", i%labels)
	}

	asStrings := heapAfter(func() any {
		out := make([]string, 0, records)
		for i := range records {
			out = append(out, makeLabel(i))
		}
		return out
	})

	asHandles := heapAfter(func() any {
		out := make([]Tag, 0, records)
		for i := range records {
			out = append(out, Intern(makeLabel(i)))
		}
		return out
	})

	Sink = nil

	if asHandles >= asStrings/2 {
		t.Errorf("interning retained %d bytes against %d for plain strings — expected less than half; is Intern returning a fresh handle each call?",
			asHandles, asStrings)
	}
}
