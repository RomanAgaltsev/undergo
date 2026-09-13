package set

import (
	"context"
	"errors"
	"strconv"
	"testing"
)

func TestMapPreservesOrder(t *testing.T) {
	s := New(1, 2, 3)

	got, err := s.Map(context.Background(), func(_ context.Context, v int) (string, error) {
		return strconv.Itoa(v), nil
	})
	if err != nil {
		t.Fatalf("Map: %v", err)
	}

	want := []string{"1", "2", "3"}
	if len(got) != len(want) {
		t.Fatalf("got %d results, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("result %d = %q, want %q", i, got[i], want[i])
		}
	}
}

// TestMapInfersR is a compile-time assertion as much as a run-time one: the call
// site writes no type argument, and the result is assigned to a declared []bool.
// A wrong inference fails the build rather than the test.
func TestMapInfersR(t *testing.T) {
	s := New("a", "bb", "ccc")

	var got []bool
	got, err := s.Map(context.Background(), func(_ context.Context, v string) (bool, error) {
		return len(v) > 1, nil
	})
	if err != nil {
		t.Fatalf("Map: %v", err)
	}
	if len(got) != 3 || got[0] || !got[1] || !got[2] {
		t.Errorf("got %v, want [false true true]", got)
	}
}

func TestEmptySetNeverCallsFn(t *testing.T) {
	s := New[int]()

	calls := 0
	got, err := s.Map(context.Background(), func(_ context.Context, v int) (int, error) {
		calls++
		return v, nil
	})
	if err != nil {
		t.Fatalf("Map: %v", err)
	}
	if got != nil {
		t.Errorf("got %v, want a nil slice", got)
	}
	if calls != 0 {
		t.Errorf("fn was called %d times, want 0", calls)
	}
}

var errBoom = errors.New("boom")

func TestFirstErrorStopsAndReportsIndex(t *testing.T) {
	s := New(0, 1, 2, 3, 4)

	calls := 0
	_, err := s.Map(context.Background(), func(_ context.Context, v int) (int, error) {
		calls++
		if v == 2 {
			return 0, errBoom
		}
		return v, nil
	})
	if err == nil {
		t.Fatal("Map returned no error")
	}

	var mapErr *MapError
	if !errors.As(err, &mapErr) {
		t.Fatalf("error %v is not a *MapError", err)
	}
	if mapErr.Index != 2 {
		t.Errorf("MapError.Index = %d, want 2", mapErr.Index)
	}
	if !errors.Is(err, errBoom) {
		t.Error("errors.Is could not reach the underlying error: is Unwrap implemented?")
	}
	if calls != 3 {
		t.Errorf("fn was called %d times, want 3 — the walk must stop at the failure", calls)
	}
}

func TestCancelledContextIsCheckedFirst(t *testing.T) {
	s := New(1, 2, 3)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	calls := 0
	_, err := s.Map(ctx, func(_ context.Context, v int) (int, error) {
		calls++
		return v, nil
	})
	if !errors.Is(err, context.Canceled) {
		t.Errorf("err = %v, want context.Canceled", err)
	}
	if calls != 0 {
		t.Errorf("fn was called %d times, want 0", calls)
	}
}

// TestMapAllocatesOneSlice is the point of the task: an implementation that
// appends to a nil slice passes every test above and fails this one.
func TestMapAllocatesOneSlice(t *testing.T) {
	s := New(make([]int, 64)...)
	ctx := context.Background()

	avg := testing.AllocsPerRun(100, func() {
		got, err := s.Map(ctx, func(_ context.Context, v int) (int, error) {
			return v * 2, nil
		})
		if err != nil {
			t.Fatalf("Map: %v", err)
		}
		Sink = got
	})
	if avg != 1 {
		t.Errorf("Map allocated %v times, want exactly 1 — size the result up front", avg)
	}
}
