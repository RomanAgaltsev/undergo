// Package drill — C3/05 error-not-returned.
package drill

import "fmt"

// Validate checks name and age, returning an error describing the first problem.
func Validate(name string, age int) error {
	if name == "" {
		fmt.Errorf("name is required")
	}
	if age < 0 {
		return fmt.Errorf("age must be non-negative, got %d", age)
	}
	return nil
}
