package cache

import (
	"runtime"
	"sync"
	"testing"
	"time"
)

// Payload is big enough to be an ordinary heap object.
type Payload struct{ Data [1024]byte }

// settle is how long a test waits for cleanups. They run on another goroutine
// some time after the collection that queued them, and CI runners are loaded,
// so this is generous: a correct cache settles in milliseconds.
const settle = 5 * time.Second

// eventually collects and polls cond until it holds or settle runs out.
func eventually(cond func() bool) bool {
	deadline := time.Now().Add(settle)
	for time.Now().Before(deadline) {
		runtime.GC()
		if cond() {
			return true
		}
		time.Sleep(time.Millisecond)
	}
	return false
}

func TestGetReturnsAValueThatIsStillHeld(t *testing.T) {
	c := New[string, Payload]()
	v := &Payload{}
	c.Put("a", v)

	runtime.GC()

	got, ok := c.Get("a")
	if !ok || got != v {
		t.Errorf("Get returned the wrong value or ok=%v while a strong reference was still held", ok)
	}
	runtime.KeepAlive(v)
}

// TestEntryRemovesItself fails for a cache that holds values strongly, for one
// with no cleanup at all, and for one whose cleanup keeps the value reachable.
func TestEntryRemovesItself(t *testing.T) {
	c := New[string, Payload]()
	func() { c.Put("a", &Payload{}) }()

	if !eventually(func() bool { return c.Len() == 0 }) {
		t.Fatalf("Len = %d %v after the value became unreachable: the entry never removed itself", c.Len(), settle)
	}
	if _, ok := c.Get("a"); ok {
		t.Error("Get found an entry that had removed itself")
	}
}

// TestReplacedEntrySurvivesTheOldCleanup is the test this task is about: the
// first value's cleanup must not delete the entry a later Put made.
func TestReplacedEntrySurvivesTheOldCleanup(t *testing.T) {
	c := New[string, Payload]()
	func() { c.Put("a", &Payload{}) }()
	v2 := &Payload{}
	c.Put("a", v2)
	func() { c.Put("sentinel", &Payload{}) }()

	// The sentinel dies in the same collection as the first value. Cleanups run
	// in no specified order, so once the sentinel's entry is gone keep watching
	// a while longer for the first value's cleanup to do damage.
	var goneAt time.Time
	settled := eventually(func() bool {
		if got, ok := c.Get("a"); !ok || got != v2 {
			t.Fatalf("the entry Put for the second value was removed — did the first value's cleanup delete it?")
		}
		if goneAt.IsZero() {
			if c.Len() == 1 {
				goneAt = time.Now()
			}
			return false
		}
		return time.Since(goneAt) > 200*time.Millisecond
	})
	if !settled {
		t.Fatalf("the sentinel entry never removed itself (Len = %d)", c.Len())
	}
	runtime.KeepAlive(v2)
}

// TestPutNilRemoves pins the contract for a nil value. runtime.AddCleanup
// panics on a nil pointer, so a Put that does not handle nil fails here.
func TestPutNilRemoves(t *testing.T) {
	c := New[string, Payload]()
	v := &Payload{}
	c.Put("a", v)
	c.Put("a", nil)

	if _, ok := c.Get("a"); ok {
		t.Error("Get found a key after Put(k, nil)")
	}
	if got := c.Len(); got != 0 {
		t.Errorf("Len = %d after Put(k, nil), want 0", got)
	}
	runtime.KeepAlive(v)
}

// liveHeap collects and returns the bytes the heap still holds.
func liveHeap() uint64 {
	runtime.GC()
	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.HeapAlloc
}

// TestRePutDoesNotAccumulate refreshes one key with the same live value many
// times, as a cache does. Each Put registers a cleanup; one that never undoes
// the previous registration keeps every earlier one — and whatever it captured
// — alive for as long as the value lives. A correct cache holds one per entry.
func TestRePutDoesNotAccumulate(t *testing.T) {
	const puts = 200_000
	const allowed = 1 << 20

	c := New[string, Payload]()
	v := &Payload{}
	c.Put("a", v)
	before := liveHeap()
	for range puts {
		c.Put("a", v)
	}
	grown := int64(liveHeap()) - int64(before)
	if grown > allowed {
		t.Errorf("the live heap grew by %d bytes over %d re-Puts of one live value — does Put leave the previous cleanup registered?", grown, puts)
	}
	runtime.KeepAlive(v)
}

// TestConcurrentUse gives the race detector something to find: Put, Get and
// Len from several goroutines while earlier values are collected and their
// cleanups run on goroutines of their own. Without -race it can only trip the
// runtime's "concurrent map writes" check, which is not guaranteed to fire.
func TestConcurrentUse(t *testing.T) {
	c := New[int, Payload]()
	var wg sync.WaitGroup
	for g := range 4 {
		wg.Go(func() {
			for i := range 5000 {
				c.Put((g*5000+i)%64, &Payload{})
				c.Get(i % 64)
				c.Len()
			}
		})
	}
	wg.Wait()

	if !eventually(func() bool { return c.Len() == 0 }) {
		t.Errorf("Len = %d %v after every value became unreachable", c.Len(), settle)
	}
}
