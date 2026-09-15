// Package edges holds eight two-goroutine programs that differ in exactly one
// thing: the act that separates the write from the read.
//
// In every one of them a goroutine sets payload to 42 and the main goroutine
// then reads it. Each function reports whether the read saw 42 on this run.
//
// The question is not what they return. It is which of them are *guaranteed*
// to return true, and which merely do.
package edges

import (
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// UnbufferedSend writes, then sends on an unbuffered channel. The reader
// receives before reading.
func UnbufferedSend() bool {
	var payload int
	ch := make(chan struct{})
	go func() {
		payload = 42
		ch <- struct{}{}
	}()
	<-ch
	return payload == 42
}

// BufferedSendFull writes, then sends on a channel of capacity 1 that is
// already full. The reader drains one value first, then receives.
func BufferedSendFull() bool {
	var payload int
	ch := make(chan struct{}, 1)
	ch <- struct{}{} // fill it
	go func() {
		payload = 42
		ch <- struct{}{} // blocks until the reader takes one out
	}()
	<-ch // frees the slot
	<-ch // receives the writer's send
	return payload == 42
}

// CloseChannel writes, then closes. The reader receives the zero value.
func CloseChannel() bool {
	var payload int
	ch := make(chan struct{})
	go func() {
		payload = 42
		close(ch)
	}()
	<-ch
	return payload == 42
}

// MutexUnlockLock writes under a lock the reader then takes.
func MutexUnlockLock() bool {
	var payload int
	var mu sync.Mutex
	mu.Lock()
	go func() {
		payload = 42
		mu.Unlock()
	}()
	mu.Lock()
	defer mu.Unlock()
	return payload == 42
}

// OnceDo has both goroutines call Do on the same Once. Only one runs the
// function; the other's Do returns after it has.
func OnceDo() bool {
	var payload int
	var once sync.Once
	var started sync.WaitGroup
	started.Add(1)

	go func() {
		started.Done()
		once.Do(func() { payload = 42 })
	}()
	started.Wait()
	once.Do(func() { payload = 42 })
	return payload == 42
}

// WaitGroupWait writes, then calls Done. The reader returns from Wait.
func WaitGroupWait() bool {
	var payload int
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		payload = 42
		wg.Done()
	}()
	wg.Wait()
	return payload == 42
}

// AtomicStoreLoad writes, then stores a flag. The reader spins on the flag.
func AtomicStoreLoad() bool {
	var payload int
	var ready atomic.Bool
	go func() {
		payload = 42
		ready.Store(true)
	}()
	for !ready.Load() {
		runtime.Gosched()
	}
	return payload == 42
}

// SleepOnly writes with no synchronisation at all. The reader waits a
// millisecond, which is not a synchronising act — it is only a long time.
func SleepOnly() bool {
	var payload int
	go func() { payload = 42 }()
	time.Sleep(time.Millisecond)
	return payload == 42
}

// NoSync is SleepOnly without the sleep. It has exactly the same
// synchronisation as SleepOnly — none — and differs only in how long the reader
// waits before looking.
func NoSync() bool {
	var payload int
	go func() { payload = 42 }()
	return payload == 42
}
