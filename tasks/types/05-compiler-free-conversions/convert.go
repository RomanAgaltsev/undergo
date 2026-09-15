// Package convert asks which string([]byte) conversions actually copy.
//
// Converting a []byte to a string normally allocates: strings are immutable and
// slices are not, so the bytes must be copied or a later write would change the
// string. The compiler skips the copy in the cases where it can prove the
// string does not outlive the expression.
//
// Every sink is typed. An `any` sink would box the value and add an allocation
// of its own, which would be measuring the test rather than the subject.
package convert

import "testing"

// Sinks keep the results from being optimised away. They are exported because
// a package-level variable that is only ever written is flagged as unused.
var (
	SinkInt    int
	SinkBool   bool
	SinkString string
)

var lookup = map[string]int{"hello": 1}

// MapIndex indexes a map with a converted key.
func MapIndex(b []byte) float64 {
	return testing.AllocsPerRun(1000, func() { SinkInt = lookup[string(b)] })
}

// Comparison compares a converted slice against a constant.
func Comparison(b []byte) float64 {
	return testing.AllocsPerRun(1000, func() { SinkBool = string(b) == "hello" })
}

// RangeOver ranges over a converted slice.
func RangeOver(b []byte) float64 {
	return testing.AllocsPerRun(1000, func() {
		n := 0
		for range string(b) {
			n++
		}
		SinkInt = n
	})
}

// TypeSwitch switches on a converted slice.
func TypeSwitch(b []byte) float64 {
	return testing.AllocsPerRun(1000, func() {
		switch string(b) {
		case "hello":
			SinkBool = true
		default:
			SinkBool = false
		}
	})
}

// Assignment assigns the converted slice to a variable that outlives the
// expression.
func Assignment(b []byte) float64 {
	return testing.AllocsPerRun(1000, func() { SinkString = string(b) })
}
