package cache

import (
	"runtime"
	"testing"
)

// Payload is big enough to be an ordinary heap object.
type Payload struct{ Data [1024]byte }

func TestGetReturnsAValueThatIsStillHeld(t *testing.T) {
	c := New[string, Payload]()
	v := &Payload{}
	c.Put("a", v)

	runtime.GC()

	got, ok := c.Get("a")
	if !ok || got != v {
		t.Errorf("Get returned the wrong value or ok=%v while a strong reference was still held", ok)
	}
	// Without this the compiler is free to consider v dead before the GC above,
	// and the test would be asserting the opposite of what it means to.
	runtime.KeepAlive(v)
}

// TestValueIsDroppedOnceNothingElseHoldsIt is the test a strong cache fails.
//
// The allocation happens inside a function literal so that no reference to it
// survives on this test's stack frame.
func TestValueIsDroppedOnceNothingElseHoldsIt(t *testing.T) {
	c := New[string, Payload]()
	func() {
		c.Put("a", &Payload{})
	}()

	runtime.GC()

	if _, ok := c.Get("a"); ok {
		t.Error("Get still returned a value after it became unreachable — is the cache holding it strongly?")
	}
}

func TestLenCountsEntriesNotLiveValues(t *testing.T) {
	c := New[string, Payload]()
	func() {
		c.Put("a", &Payload{})
	}()

	runtime.GC()

	if got := c.Len(); got != 1 {
		t.Errorf("Len = %d, want 1 — the entry outlives the value until it is reaped", got)
	}
	c.Reap()
	if got := c.Len(); got != 0 {
		t.Errorf("Len = %d after Reap, want 0", got)
	}
}

func TestReapKeepsLiveEntries(t *testing.T) {
	c := New[string, Payload]()
	alive := &Payload{}
	c.Put("alive", alive)
	func() {
		c.Put("dead", &Payload{})
	}()

	runtime.GC()
	c.Reap()

	if got := c.Len(); got != 1 {
		t.Errorf("Len = %d after Reap, want 1 — Reap must keep entries whose value is still held", got)
	}
	if _, ok := c.Get("alive"); !ok {
		t.Error("Reap removed an entry whose value is still reachable")
	}
	runtime.KeepAlive(alive)
}

func TestPutReplacesAnExistingKey(t *testing.T) {
	c := New[string, Payload]()
	first, second := &Payload{}, &Payload{}
	c.Put("a", first)
	c.Put("a", second)

	if got, ok := c.Get("a"); !ok || got != second {
		t.Errorf("Get returned the wrong value or ok=%v, want the second value", ok)
	}
	if got := c.Len(); got != 1 {
		t.Errorf("Len = %d after replacing a key, want 1", got)
	}
	runtime.KeepAlive(first)
	runtime.KeepAlive(second)
}

func TestMissingKey(t *testing.T) {
	c := New[string, Payload]()
	if _, ok := c.Get("nope"); ok {
		t.Error("Get reported ok for a key that was never stored")
	}
}
