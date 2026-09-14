// Package dispatch holds five interface call sites.
//
// You do not implement anything here. Work out which the compiler resolves
// statically, write your answers into prediction.yaml, then verify.
package dispatch

// Shape is the interface every call below goes through.
type Shape interface{ Area() int }

// Sq is one implementation.
type Sq struct{ S int }

// Area satisfies Shape.
func (s Sq) Area() int { return s.S * s.S }

// Rect is a second implementation, so the interface has more than one.
type Rect struct{ W, H int }

// Area satisfies Shape.
func (r Rect) Area() int { return r.W * r.H }

// Global holds a Shape assigned elsewhere.
var Global Shape = Sq{S: 2}

// OneType assigns a single concrete type and calls through the interface.
func OneType() int {
	var s Shape = Sq{S: 4}
	return s.Area()
}

// Reassigned assigns one type and then another before calling.
func Reassigned(pick bool) int {
	var s Shape = Sq{S: 4}
	if pick {
		s = Rect{W: 2, H: 3}
	}
	return s.Area()
}

// ViaParameter is handed the interface by its caller.
func ViaParameter(s Shape) int { return s.Area() }

// ViaGlobal calls through a package-level interface variable.
func ViaGlobal() int { return Global.Area() }

// newSq returns a concrete type.
func newSq() Sq { return Sq{S: 4} }

// FromConstructor assigns the result of a function returning a concrete type.
func FromConstructor() int {
	var s Shape = newSq()
	return s.Area()
}

// InSlice calls through an element of a slice of interfaces.
func InSlice() int {
	shapes := []Shape{Sq{S: 2}}
	t := 0
	for _, s := range shapes {
		t += s.Area()
	}
	return t
}
