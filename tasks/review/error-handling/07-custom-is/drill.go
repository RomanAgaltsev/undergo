// Package drill — C3/07 custom-is.
package drill

import (
	"errors"
	"fmt"
)

// NotFoundError indicates a missing key.
type NotFoundError struct {
	Key string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("key %q not found", e.Key)
}

// ErrNotFound is the sentinel some callers check for.
var ErrNotFound = errors.New("not found")

// Lookup always fails with a wrapped NotFoundError for the given key.
func Lookup(key string) error {
	return fmt.Errorf("lookup: %w", &NotFoundError{Key: key})
}

// IsMissing reports whether err denotes a missing key.
func IsMissing(err error) bool {
	return errors.Is(err, ErrNotFound)
}
