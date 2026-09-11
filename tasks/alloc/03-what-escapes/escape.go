// Package escape holds four functions, each creating one Holder. Some of those
// Holders escape to the heap and some do not.
//
// You do not implement anything here. Work out which is which, write your
// answers into prediction.yaml, then verify — the compiler is the judge.
package escape

// Holder is the value under test in each function below.
type Holder struct{ N int }

// Sink keeps escaping values reachable.
var Sink *Holder

// Local builds a Holder and reads one field out of it.
func Local(n int) int {
	h := Holder{N: n}
	return h.N
}

// Returned builds a Holder and returns its address.
func Returned(n int) *Holder {
	h := Holder{N: n}
	return &h
}

// Stored builds a Holder and parks its address in a package variable.
func Stored(n int) {
	h := Holder{N: n}
	Sink = &h
}

// Interfaced builds a Holder and returns it inside an any. Nothing takes its
// address.
func Interfaced(n int) any {
	h := Holder{N: n}
	return any(h)
}
