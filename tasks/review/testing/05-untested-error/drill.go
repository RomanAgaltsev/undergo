// Package drill — C10/05 untested-error.
package drill

import "errors"

// Withdraw subtracts amount from balance, erroring on overdraft.
func Withdraw(balance, amount int) (int, error) {
	if amount > balance {
		return balance, errors.New("insufficient funds")
	}
	return balance - amount, nil
}
