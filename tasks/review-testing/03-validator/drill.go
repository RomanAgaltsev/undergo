// Package drill — C10/03 validator.
package drill

import "errors"

// Validator validates names.
type Validator struct{}

// Check reports whether name is valid: it must be non-empty (else an error) and
// at most 10 characters (else false).
func (Validator) Check(name string) (bool, error) {
	if name == "" {
		return false, errors.New("empty name")
	}
	if len(name) > 10 {
		return false, nil
	}
	return true, nil
}
