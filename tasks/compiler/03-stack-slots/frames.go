// Package frames holds five functions whose stack frames you are asked to
// compare.
//
// You do not implement anything here. Work out the answers, write them into
// prediction.yaml, then verify.
//
// Every slot is a comparison rather than a byte count, because a frame size is
// not portable: the same function measures eight bytes larger on arm64 than on
// amd64. Ratios survive that; absolute sizes do not.
package frames

// Sink keeps results reachable so nothing here is optimised away.
var Sink int

// use takes the address of an array so it cannot be kept in registers.
//
//go:noinline
func use(p *[1024]byte, n int) int {
	p[0] = byte(n)
	return int(p[0])
}

// Helper is a leaf that needs no storage of its own.
func Helper(n int) int { return n * 2 }

// Unused declares a large array and never does anything the compiler cannot
// see through.
func Unused(n int) int {
	var a [1024]byte
	return n + int(a[0])
}

// One keeps a single large array alive.
func One(n int) int {
	var a [1024]byte
	return use(&a, n)
}

// TwoLive keeps two large arrays alive at the same time.
func TwoLive(n int) int {
	var a, b [1024]byte
	t := use(&a, n) + use(&b, n)
	return t + int(a[0]) + int(b[0])
}

// TwoDisjoint uses two large arrays in separate scopes, so their live ranges do
// not overlap.
func TwoDisjoint(n int) int {
	t := 0
	{
		var a [1024]byte
		t += use(&a, n)
	}
	{
		var b [1024]byte
		t += use(&b, n)
	}
	return t
}
