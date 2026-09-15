//go:build amd64

// Package add128 adds two 128-bit numbers held as pairs of uint64.
//
// Go has no way to read the carry flag. math/bits.Add64 exists precisely
// because the language cannot express what one instruction does. Write that
// instruction and find out what the compiler has been doing for you.
package add128

// Add128 returns the 128-bit sum of (hiA:loA) and (hiB:loB), discarding any
// carry out of the top.
//
// Implement it in add128_amd64.s. The stub there adds the low words and throws
// the carry away, which is exactly the bug you are here to fix.
func Add128(loA, hiA, loB, hiB uint64) (lo, hi uint64)
