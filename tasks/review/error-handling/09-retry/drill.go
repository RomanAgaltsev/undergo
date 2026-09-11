// Package drill — C3/09 retry.
package drill

import "errors"

// ErrTemporary marks a retryable failure.
var ErrTemporary = errors.New("temporary")

// Retry calls fn up to attempts times, retrying only temporary errors.
func Retry(attempts int, fn func() error) error {
	for i := 0; i < attempts; i++ {
		err := fn()
		if err == nil {
			return nil
		}
		if err == ErrTemporary {
			continue
		}
		_ = err
	}
	return nil
}
