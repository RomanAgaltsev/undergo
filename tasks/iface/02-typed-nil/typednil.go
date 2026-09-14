// Package typednil is about an interface value holding a nil pointer.
//
// You do not implement anything here. Work out the answers, write them into
// prediction.yaml, then verify.
package typednil

import (
	"errors"
	"reflect"
)

// MyErr is a concrete error type.
type MyErr struct{ Code int }

// Error satisfies the error interface.
func (e *MyErr) Error() string { return "my error" }

// mayFail returns a nil *MyErr — note the concrete return type.
func mayFail() *MyErr { return nil }

// TrueNilIsNil reports whether an interface that was never assigned is nil.
func TrueNilIsNil() bool {
	var err error
	return err == nil
}

// NilPointerInInterfaceIsNil reports whether an interface holding a nil pointer
// compares equal to nil.
func NilPointerInInterfaceIsNil() bool {
	var p *MyErr //nolint:staticcheck // the whole point is that this concrete nil is visible to SA4023
	var err error = p
	//nolint:staticcheck // SA4023 is right, and being right is the point of this task.
	return err == nil
}

// ReturnedNilPointerIsNil is the same thing arriving from a function.
func ReturnedNilPointerIsNil() bool {
	var err error = mayFail() //nolint:staticcheck // deliberately returns a concrete nil
	//nolint:staticcheck // SA4023 is right, and being right is the point of this task.
	return err == nil
}

// ComparingToTypedNilWorks reports whether comparing against a typed nil finds
// the value.
func ComparingToTypedNilWorks() bool {
	var p *MyErr //nolint:staticcheck // the whole point is that this concrete nil is visible to SA4023
	var err error = p
	//nolint:errorlint // comparing against a typed nil is the question being asked.
	return err == (*MyErr)(nil)
}

// ReflectSaysNil reports what reflect thinks of the same value.
func ReflectSaysNil() bool {
	var p *MyErr //nolint:staticcheck // the whole point is that this concrete nil is visible to SA4023
	var err error = p
	return reflect.ValueOf(err).IsNil()
}

// ErrorsIsNilAgreesWithEquality reports whether errors.Is(err, nil) gives the
// same answer as err == nil for this value.
func ErrorsIsNilAgreesWithEquality() bool {
	var p *MyErr //nolint:staticcheck // the whole point is that this concrete nil is visible to SA4023
	var err error = p
	//nolint:staticcheck // SA4023 again, deliberately.
	return errors.Is(err, nil) == (err == nil)
}
