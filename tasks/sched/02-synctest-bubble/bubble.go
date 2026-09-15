// Package bubble exercises testing/synctest's fake clock.
//
// Inside a bubble the time package runs on a per-bubble clock that starts at
// midnight UTC on 2000-01-01 and advances only when every goroutine in the
// bubble is durably blocked.
package bubble

import "time"

// SleepTwice sleeps for a, then b, and reports how much time the time package
// believes has passed.
func SleepTwice(a, b time.Duration) time.Duration {
	start := time.Now()
	time.Sleep(a)
	time.Sleep(b)
	return time.Since(start)
}

// TickCount counts how many times a ticker of period p fires while the caller
// waits for total. The ticker is stopped before returning.
func TickCount(p, total time.Duration) int {
	tk := time.NewTicker(p)
	defer tk.Stop()

	done := time.After(total)
	n := 0
	for {
		select {
		case <-tk.C:
			n++
		case <-done:
			return n
		}
	}
}

// Year reports the year the time package believes it is.
func Year() int { return time.Now().Year() }
