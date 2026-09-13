package stamp

import (
	"testing"
	"time"
)

var when = time.Date(2026, 9, 13, 12, 0, 0, 0, time.UTC)

func TestStampsInPlace(t *testing.T) {
	events := []Event{{Name: "a"}, {Name: "b"}}

	StampAll[Event](events, when)

	for i := range events {
		if !events[i].Stamped.Equal(when) {
			t.Errorf("events[%d] (%q) was not stamped — did the loop range by value?", i, events[i].Name)
		}
	}
}

func TestReturnsTheSameBackingArray(t *testing.T) {
	events := []Event{{Name: "a"}, {Name: "b"}}

	got := StampAll[Event](events, when)

	if len(got) != len(events) {
		t.Fatalf("got %d items, want %d", len(got), len(events))
	}
	if &got[0] != &events[0] {
		t.Error("the result does not share a backing array with the argument")
	}
}

func TestEmptyAndNil(t *testing.T) {
	if got := StampAll[Event](nil, when); len(got) != 0 {
		t.Errorf("nil input returned %d items", len(got))
	}
	if got := StampAll[Event]([]Event{}, when); len(got) != 0 {
		t.Errorf("empty input returned %d items", len(got))
	}
}

// TestDoesNotAllocate catches an implementation that sidesteps the constraint by
// building a new slice and stamping that.
//
// It does NOT catch converting each element to a Stampable — an interface
// holding a pointer needs no boxing allocation, so that version measures zero
// too. Written question 4 asks you to work out why.
func TestDoesNotAllocate(t *testing.T) {
	events := make([]Event, 64)

	avg := testing.AllocsPerRun(100, func() {
		StampAll[Event](events, when)
	})
	if avg != 0 {
		t.Errorf("StampAll allocated %v times, want 0 — is a slice being copied?", avg)
	}
}
