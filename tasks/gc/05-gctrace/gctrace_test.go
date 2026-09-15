package gctrace

import (
	"os"
	"strings"
	"testing"
)

func fixture(t *testing.T) string {
	t.Helper()
	b, err := os.ReadFile("testdata/trace.txt")
	if err != nil {
		t.Fatalf("reading the fixture: %v", err)
	}
	return string(b)
}

// TestParsesEveryCollection checks the count and the ends.
func TestParsesEveryCollection(t *testing.T) {
	cs, err := Parse(fixture(t))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(cs) != 15 {
		t.Fatalf("expected 15 collections, got %d", len(cs))
	}
	if cs[0].Num != 1 {
		t.Errorf("the first collection is numbered wrong")
	}
	if cs[len(cs)-1].Num != 15 {
		t.Errorf("the last collection is numbered wrong")
	}
}

// TestIgnoresNonCollectionLines checks that the scavenger line and the bare
// "GC forced" banner do not become cycles. Note the banner contains the word
// "forced", which a careless parser will notice.
func TestIgnoresNonCollectionLines(t *testing.T) {
	trace := fixture(t)
	if !strings.Contains(trace, "scvg:") || !strings.Contains(trace, "GC forced") {
		t.Fatal("the fixture no longer contains the lines this test is about")
	}
	cs, err := Parse(trace)
	if err != nil {
		t.Fatalf("Parse rejected a trace containing non-collection lines: %v", err)
	}
	if len(cs) != 15 {
		t.Errorf("a non-collection line was parsed as a collection")
	}
}

// TestFirstCycleFields checks every field of a known line:
//
//	gc 1 @0.001s 4%: 0+0.52+0 ms clock, ... 3->4->0 MB, 4 MB goal, ... 12 P
func TestFirstCycleFields(t *testing.T) {
	cs, err := Parse(fixture(t))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	c := cs[0]
	switch {
	case c.Num != 1:
		t.Error("Num is wrong")
	case c.PercentCPU != 4:
		t.Error("PercentCPU is wrong")
	case c.HeapStart != 3:
		t.Error("HeapStart is wrong")
	case c.HeapEnd != 0:
		t.Error("HeapEnd is wrong")
	case c.HeapGoal != 4:
		t.Error("HeapGoal is wrong")
	case c.Procs != 12:
		t.Error("Procs is wrong")
	case c.Forced:
		t.Error("the first collection was not forced")
	case c.WallMS < 0.51 || c.WallMS > 0.53:
		t.Errorf("WallMS should be the three clock figures added together")
	}
}

// TestForcedCollections checks the flag and the count.
func TestForcedCollections(t *testing.T) {
	cs, err := Parse(fixture(t))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if got := TotalForced(cs); got != 2 {
		t.Errorf("expected 2 forced collections, got %d", got)
	}
	forced := map[int]bool{}
	for _, c := range cs {
		if c.Forced {
			forced[c.Num] = true
		}
	}
	if !forced[4] || !forced[11] {
		t.Errorf("the wrong collections are marked forced")
	}
}

// TestMalformedLineIsAnError checks that a broken collection line is reported
// rather than silently yielding a zero Cycle.
func TestMalformedLineIsAnError(t *testing.T) {
	if _, err := Parse("gc 1 @0.001s this is not a trace line\n"); err == nil {
		t.Error("a malformed collection line should be an error")
	}
}

// TestEmptyTrace checks the empty case.
func TestEmptyTrace(t *testing.T) {
	cs, err := Parse("")
	if err != nil {
		t.Fatalf("an empty trace is not an error: %v", err)
	}
	if len(cs) != 0 {
		t.Error("an empty trace has no collections")
	}
}
