// Package spinlock builds a lock out of one atomic value.
//
// A sync.Mutex parks a goroutine that cannot get in. A spin lock never parks:
// a goroutine that finds it taken tries again, and again, until it gets it.
// That is cheap when the wait is short and the waiter has a core of its own,
// and a disaster when neither is true.
package spinlock

import "sync"

var _ sync.Locker = (*SpinLock)(nil)

// SpinLock is a mutual-exclusion lock that never parks a goroutine. The zero
// value is unlocked. It must not be copied after first use.
type SpinLock struct {
	// Declare what you need. The contract asks for one value from sync/atomic
	// and nothing from sync.
}

// Lock takes the lock, retrying until it succeeds. While it retries it must
// let other goroutines run.
func (l *SpinLock) Lock() {
	panic("undergo: implement Lock")
}

// Unlock releases the lock. Unlocking a SpinLock that is not locked is a bug
// in the caller and need not be detected.
func (l *SpinLock) Unlock() {
	panic("undergo: implement Unlock")
}
