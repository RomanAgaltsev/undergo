// Package drill — C15/04 map-of-iface.
package drill

// fieldError is a per-field validation failure.
type fieldError struct {
	field string
}

func (e *fieldError) Error() string { return "invalid: " + e.field }

// Form is the data being validated.
type Form struct {
	Name  string
	Email string
}

// checkName returns a non-nil error if the name is missing.
func checkName(f Form) error {
	var e *fieldError
	if f.Name == "" {
		e = &fieldError{field: "name"}
	}
	return e
}

// checkEmail returns a non-nil error if the email is missing.
func checkEmail(f Form) error {
	var e *fieldError
	if f.Email == "" {
		e = &fieldError{field: "email"}
	}
	return e
}

// Problems returns the list of validation errors for f (empty if valid).
func Problems(f Form) []error {
	checks := []func(Form) error{checkName, checkEmail}

	var problems []error
	for _, check := range checks {
		if err := check(f); err != nil {
			problems = append(problems, err)
		}
	}
	return problems
}
