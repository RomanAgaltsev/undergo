// Package timers looks at two things Go changed about timers: the buffering of
// the channel a Timer or Ticker hands you, and whether one you drop without
// calling Stop can be collected.
package timers

import (
	"runtime"
	"time"
)

// TimerCap reports the capacity of the channel inside a Timer.
func TimerCap() int {
	t := time.NewTimer(time.Hour)
	defer t.Stop()
	return cap(t.C)
}

// TickerCap reports the capacity of the channel inside a Ticker.
func TickerCap() int {
	t := time.NewTicker(time.Hour)
	defer t.Stop()
	return cap(t.C)
}

// AfterCap reports the capacity of the channel time.After returns.
func AfterCap() int { return cap(time.After(time.Hour)) }

// DroppedTickerCollected creates a Ticker, never calls Stop, drops the last
// reference to it, and reports whether the garbage collector reclaimed it.
//
// The answer is not about the collector's cleverness. It is about whether the
// runtime still holds a reference to a running ticker.
func DroppedTickerCollected() bool {
	collected := make(chan struct{})

	func() {
		t := time.NewTicker(time.Hour)
		runtime.AddCleanup(t, func(ch chan struct{}) { close(ch) }, collected)
	}()

	for range 10 {
		runtime.GC()
		select {
		case <-collected:
			return true
		default:
		}
	}
	return false
}
