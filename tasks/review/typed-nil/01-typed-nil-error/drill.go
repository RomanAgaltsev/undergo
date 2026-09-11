// Package drill — C15/01 typed-nil-error.
package drill

import "fmt"

// Input is the value being validated.
type Input struct {
	Name string
	Age  int
}

// ValidationError describes why an input was rejected.
type ValidationError struct {
	Field string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("invalid field: %s", e.Field)
}

// Validate checks in and returns a non-nil error if it is invalid.
func Validate(in Input) error {
	var e *ValidationError
	if in.Name == "" {
		e = &ValidationError{Field: "name"}
	}
	return e
}
