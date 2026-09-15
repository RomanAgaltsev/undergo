//go:build amd64

// Package noescape asks what //go:noescape is worth.
//
// Both functions below are implemented in noescape_amd64.s and have byte-
// identical bodies: they sum four int64s behind a pointer. Only one of them
// carries the pragma.
//
// The compiler cannot see into an assembly file, so it has to make an
// assumption about every pointer that crosses into one. The pragma replaces
// the assumption with your word.
package noescape

// SumFree sums the four elements behind p. Implemented in assembly, and
// declared with the pragma that says the assembly does not retain p.
//
//go:noescape
func SumFree(p *[4]int64) int64

// SumKept sums the four elements behind p. The assembly is byte-identical to
// SumFree's; only the pragma is missing.
func SumKept(p *[4]int64) int64

// CallFree builds an array on the stack, if it can, and sums it through the
// function that carries the pragma.
func CallFree() int64 {
	var a [4]int64
	a[0], a[1], a[2], a[3] = 1, 2, 3, 4
	return SumFree(&a)
}

// CallKept does the same through the function that does not.
func CallKept() int64 {
	var a [4]int64
	a[0], a[1], a[2], a[3] = 1, 2, 3, 4
	return SumKept(&a)
}
