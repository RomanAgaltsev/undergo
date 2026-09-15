// Package deepequal is a structural equality check that survives cycles.
//
// reflect.DeepEqual already does this. Write your own and find out what it
// costs: the naive recursive version does not terminate on a value that
// contains itself, and the fix is not the one people reach for first.
package deepequal

// Equal reports whether a and b are structurally equal.
//
// It must follow pointers, slices, maps and structs, compare their contents
// rather than their addresses, and terminate even when either value contains a
// cycle.
//
// Match reflect.DeepEqual's semantics: a nil slice and an empty slice are not
// equal, and neither are a nil map and an empty map.
func Equal(a, b any) bool {
	panic("implement Equal")
}
