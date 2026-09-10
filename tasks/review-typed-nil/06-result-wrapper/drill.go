// Package drill — C15/06 result-wrapper.
package drill

// Output is the produced value.
type Output struct {
	Data string
}

// runError describes a failed run.
type runError struct {
	reason string
}

func (e *runError) Error() string { return e.reason }

// Result wraps the outcome of a run.
type Result struct {
	Value any
	Err   error
}

// Run executes the job and wraps the outcome.
func Run(fail bool) Result {
	var out *Output
	var e *runError

	if fail {
		e = &runError{reason: "boom"}
	} else {
		out = &Output{Data: "ok"}
	}

	return Result{Value: out, Err: e}
}

// firstNonEmpty returns the first non-empty string and a nil error, or an error.
func firstNonEmpty(ss []string) (string, error) {
	var e *runError
	for _, s := range ss {
		if s != "" {
			return s, e
		}
	}
	return "", &runError{reason: "all empty"}
}
