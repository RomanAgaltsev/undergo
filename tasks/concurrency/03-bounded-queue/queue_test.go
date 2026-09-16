package boundedqueue

import (
	"errors"
	"sync"
	"testing"
	"testing/synctest"
)

func TestFIFOWithinCapacity(t *testing.T) {
	q := NewQueue[int](3)
	for i := range 3 {
		if err := q.Push(i); err != nil {
			t.Fatalf("Push(%d): %v", i, err)
		}
	}
	for want := range 3 {
		got, err := q.Pop()
		if err != nil || got != want {
			t.Fatalf("Pop() = %v, %v; want %d, nil", got, err, want)
		}
	}
}

func TestPushBlocksWhenFull(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		q := NewQueue[int](1)
		if err := q.Push(1); err != nil {
			t.Fatal(err)
		}
		pushed := make(chan struct{})
		go func() { _ = q.Push(2); close(pushed) }()
		synctest.Wait()
		select {
		case <-pushed:
			t.Fatal("Push returned while the queue was full")
		default:
		}
		if _, err := q.Pop(); err != nil {
			t.Fatal(err)
		}
		<-pushed
	})
}

func TestPopBlocksWhenEmpty(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		q := NewQueue[int](2)
		popped := make(chan int, 1)
		go func() {
			v, err := q.Pop()
			if err != nil {
				t.Error(err)
			}
			popped <- v
		}()
		synctest.Wait()
		select {
		case <-popped:
			t.Fatal("Pop returned from an empty queue")
		default:
		}
		if err := q.Push(7); err != nil {
			t.Fatal(err)
		}
		if got := <-popped; got != 7 {
			t.Fatalf("Pop() = %d; want 7", got)
		}
	})
}

// Close must wake EVERY waiter, not one of them. A Cond woken with Signal
// instead of Broadcast passes every other test in this file and hangs here.
func TestCloseWakesAllWaiters(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		q := NewQueue[int](1)
		var wg sync.WaitGroup
		for range 4 {
			wg.Go(func() {
				if _, err := q.Pop(); !errors.Is(err, ErrClosed) {
					t.Errorf("Pop() error = %v; want ErrClosed", err)
				}
			})
		}
		synctest.Wait() // all four are parked in Wait
		q.Close()
		wg.Wait()
	})
}

// Close must also wake blocked pushers, and tell them the truth.
func TestCloseWakesBlockedPushers(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		q := NewQueue[int](1)
		if err := q.Push(1); err != nil {
			t.Fatal(err)
		}
		var wg sync.WaitGroup
		for range 3 {
			wg.Go(func() {
				if err := q.Push(2); !errors.Is(err, ErrClosed) {
					t.Errorf("Push() error = %v; want ErrClosed", err)
				}
			})
		}
		synctest.Wait()
		q.Close()
		wg.Wait()
	})
}

// A closed queue still yields what it already holds. This is what forces the
// closed check into the loop condition rather than the top of Pop.
func TestCloseDrainsRemaining(t *testing.T) {
	q := NewQueue[int](4)
	if err := q.Push(1); err != nil {
		t.Fatal(err)
	}
	if err := q.Push(2); err != nil {
		t.Fatal(err)
	}
	q.Close()
	for want := 1; want <= 2; want++ {
		got, err := q.Pop()
		if err != nil || got != want {
			t.Fatalf("Pop() = %v, %v; want %d, nil", got, err, want)
		}
	}
	if _, err := q.Pop(); !errors.Is(err, ErrClosed) {
		t.Fatalf("Pop() on a drained closed queue = %v; want ErrClosed", err)
	}
}

func TestCloseIsIdempotent(t *testing.T) {
	q := NewQueue[int](1)
	q.Close()
	q.Close()
	if err := q.Push(1); !errors.Is(err, ErrClosed) {
		t.Fatalf("Push() after Close = %v; want ErrClosed", err)
	}
}

func TestConcurrentPushPop(t *testing.T) {
	const n = 200
	q := NewQueue[int](8)
	var wg sync.WaitGroup
	wg.Go(func() {
		for i := range n {
			if err := q.Push(i); err != nil {
				t.Error(err)
				return
			}
		}
		q.Close()
	})
	seen := 0
	for {
		_, err := q.Pop()
		if errors.Is(err, ErrClosed) {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		seen++
	}
	wg.Wait()
	if seen != n {
		t.Fatalf("received %d items; want %d", seen, n)
	}
}

// Many pushers against many poppers. Go's Mutex is not FIFO, so between a
// Signal and the woken goroutine reacquiring the mutex, a freshly-arriving
// caller can barge in and take the item the signal was about. A Wait guarded by
// `if` rather than `for` proceeds anyway, on a queue that is empty again.
func TestContendedPushPop(t *testing.T) {
	const (
		pushers    = 8
		poppers    = 8
		perPusher  = 250
		totalItems = pushers * perPusher
	)
	q := NewQueue[int](4)

	var produce sync.WaitGroup
	for range pushers {
		produce.Go(func() {
			for i := range perPusher {
				if err := q.Push(i); err != nil {
					t.Error(err)
					return
				}
			}
		})
	}

	var mu sync.Mutex
	counts := map[int]int{}
	var consume sync.WaitGroup
	for range poppers {
		consume.Go(func() {
			for {
				v, err := q.Pop()
				if errors.Is(err, ErrClosed) {
					return
				}
				if err != nil {
					t.Error(err)
					return
				}
				mu.Lock()
				counts[v]++
				mu.Unlock()
			}
		})
	}

	produce.Wait()
	q.Close()
	consume.Wait()

	total := 0
	for v, n := range counts {
		if v < 0 || v >= perPusher {
			t.Fatalf("received %d, which was never pushed", v)
		}
		total += n
	}
	if total != totalItems {
		t.Fatalf("received %d items; want %d", total, totalItems)
	}
	for v := range perPusher {
		if counts[v] != pushers {
			t.Fatalf("value %d received %d times; want %d", v, counts[v], pushers)
		}
	}
}

// A Broadcast wakes every waiter, but a single Push satisfies only one of them.
// The losers must go back to waiting — an open queue must never hand anybody
// ErrClosed. This is the test that a `if` in place of `for` fails: the woken
// losers fall straight through to the emptiness check and report the queue
// closed while it is still very much open.
func TestWokenLosersKeepWaiting(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		q := NewQueue[int](4)
		errs := make(chan error, 3)
		vals := make(chan int, 3)
		for range 3 {
			go func() {
				v, err := q.Pop()
				if err != nil {
					errs <- err
					return
				}
				vals <- v
			}()
		}
		synctest.Wait() // all three are parked in Wait

		if err := q.Push(99); err != nil {
			t.Fatal(err)
		}
		synctest.Wait() // the Broadcast has been seen by all three

		if len(errs) != 0 {
			t.Fatalf("a popper returned %v from an open queue", <-errs)
		}
		if len(vals) != 1 {
			t.Fatalf("%d poppers got a value; want exactly 1", len(vals))
		}
		if got := <-vals; got != 99 {
			t.Fatalf("Pop() = %d; want 99", got)
		}

		q.Close() // let the other two go
	})
}
