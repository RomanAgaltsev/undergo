// Package shapes measures what an interface value is made of.
//
// You do not implement anything here. Work out the answers, write them into
// prediction.yaml, then verify.
package shapes

import (
	"testing"
	"unsafe"
)

// Sink is typed, not any: assigning to an interface would box the value and add
// an allocation this task is not asking about.
var Sink int64

// BoxSink holds boxed values where boxing is the thing being measured.
var BoxSink any

// Speaker is an interface with one method.
type Speaker interface{ Speak() string }

// Dog implements Speaker.
type Dog struct{ Name string }

// Speak satisfies Speaker.
func (d Dog) Speak() string { return "woof" }

// EmptyInterfaceWords returns the size of an empty interface value in words.
func EmptyInterfaceWords() int {
	var v any
	return int(unsafe.Sizeof(v)) / int(unsafe.Sizeof(uintptr(0)))
}

// MethodInterfaceWords returns the size of a one-method interface value in words.
func MethodInterfaceWords() int {
	var v Speaker
	return int(unsafe.Sizeof(v)) / int(unsafe.Sizeof(uintptr(0)))
}

// SliceWords returns the size of a slice header in words, for comparison.
func SliceWords() int {
	var v []int
	return int(unsafe.Sizeof(v)) / int(unsafe.Sizeof(uintptr(0)))
}

// Small is a runtime value inside the small-integer cache.
var Small = 7

// BoxingSmallIntAllocates reports whether boxing a small int allocates.
func BoxingSmallIntAllocates() bool {
	return testing.AllocsPerRun(100, func() { BoxSink = Small }) > 0
}

// Runtime is a value the compiler cannot fold, so boxing it is a real runtime
// boxing rather than a constant hoisted into read-only data.
var Runtime = 1 << 20

// BoxingLargeIntAllocates reports whether boxing a large int allocates.
func BoxingLargeIntAllocates() bool {
	return testing.AllocsPerRun(100, func() { BoxSink = Runtime }) > 0
}

// BoxingPointerAllocates reports whether boxing a pointer allocates.
func BoxingPointerAllocates() bool {
	p := &Dog{Name: "rex"}
	return testing.AllocsPerRun(100, func() { BoxSink = p }) > 0
}

// BoxingStructAllocates reports whether boxing a struct value allocates.
func BoxingStructAllocates() bool {
	d := Dog{Name: "rex"}
	return testing.AllocsPerRun(100, func() { BoxSink = d }) > 0
}
