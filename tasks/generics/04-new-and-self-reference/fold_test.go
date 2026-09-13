package fold

import (
	"errors"
	"testing"
)

func TestSumsInOrder(t *testing.T) {
	got, err := SumAll([]Money{{100}, {250}, {5}})
	if err != nil {
		t.Fatalf("SumAll: %v", err)
	}
	if got == nil {
		t.Fatal("SumAll returned a nil pointer and no error")
	}
	if got.Cents != 355 {
		t.Errorf("total = %d, want 355", got.Cents)
	}
}

func TestSingleElement(t *testing.T) {
	got, err := SumAll([]Money{{42}})
	if err != nil {
		t.Fatalf("SumAll: %v", err)
	}
	if got == nil || got.Cents != 42 {
		t.Errorf("got %v, want a pointer to 42", got)
	}
}

func TestEmptyIsAnError(t *testing.T) {
	got, err := SumAll([]Money{})
	if !errors.Is(err, ErrEmpty) {
		t.Errorf("err = %v, want ErrEmpty", err)
	}
	if got != nil {
		t.Errorf("got %v, want a nil pointer alongside the error", got)
	}

	if _, err := SumAll[Money](nil); !errors.Is(err, ErrEmpty) {
		t.Errorf("nil slice: err = %v, want ErrEmpty", err)
	}
}

// TestResultDoesNotAliasTheInput proves the result is a new variable rather than
// the address of an element the caller still owns.
func TestResultDoesNotAliasTheInput(t *testing.T) {
	items := []Money{{1}, {2}}

	got, err := SumAll(items)
	if err != nil {
		t.Fatalf("SumAll: %v", err)
	}
	if got == &items[0] || got == &items[1] {
		t.Error("the result points into the caller's slice")
	}

	got.Cents = 999
	if items[0].Cents != 1 || items[1].Cents != 2 {
		t.Error("writing through the result changed the input")
	}
}

func TestAllocatesOnce(t *testing.T) {
	items := []Money{{1}, {2}, {3}, {4}}

	avg := testing.AllocsPerRun(100, func() {
		got, err := SumAll(items)
		if err != nil {
			t.Fatalf("SumAll: %v", err)
		}
		Sink = got
	})
	if avg != 1 {
		t.Errorf("SumAll allocated %v times, want exactly 1 — only the result should escape", avg)
	}
}
