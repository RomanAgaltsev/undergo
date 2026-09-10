// Package drill — C1/06 lock-ordering.
package drill

import "sync"

// Account holds a balance guarded by its own mutex.
type Account struct {
	mu      sync.Mutex
	balance int
}

// Transfer moves amount from a to b.
func Transfer(a, b *Account, amount int) {
	a.mu.Lock()
	defer a.mu.Unlock()
	b.mu.Lock()
	defer b.mu.Unlock()

	a.balance -= amount
	b.balance += amount
}

// Balance returns the account's balance.
func (a *Account) Balance() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.balance
}
