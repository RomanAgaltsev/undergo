// Package drill — C3/03 partial-op.
package drill

import (
	"errors"
	"fmt"
	"log"
)

// ErrConflict indicates the order already exists.
var ErrConflict = errors.New("conflict")

// Order is a customer order.
type Order struct {
	ID   int
	Paid bool
}

func writeOrder(o *Order) error { return nil }
func charge(o *Order) error     { return nil }
func notify(o *Order) error     { return nil }
func rollback(o *Order) error   { return nil }

// CreateOrder writes the order, charges the customer, then notifies downstream.
// If anything fails it is supposed to roll back and return the error.
func CreateOrder(o *Order) (err error) {
	if err = writeOrder(o); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			_ = rollback(o)
		}
	}()

	if cerr := charge(o); cerr != nil {
		log.Printf("charge failed: %v", cerr)
	}
	o.Paid = true

	if err = notify(o); err != nil {
		if err == ErrConflict {
			return nil
		}
		return fmt.Errorf("notify: %w", err)
	}
	return nil
}
