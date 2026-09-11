// Package drill — C15/03 interface-field-nil.
package drill

// opError describes a failed operation.
type opError struct {
	code int
}

func (e *opError) Error() string { return "operation failed" }

// Result carries the outcome of an operation.
type Result struct {
	Value int
	Err   error
}

// Do performs the operation and returns a Result. Err is set only on failure.
func Do(x int) Result {
	var e *opError
	if x < 0 {
		e = &opError{code: 400}
	}
	return Result{Value: x * 2, Err: e}
}

// Succeeded reports whether the result represents success.
func (r Result) Succeeded() bool {
	return r.Err == nil
}
