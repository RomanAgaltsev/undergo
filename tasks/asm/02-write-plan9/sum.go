//go:build amd64

// Package sum adds a slice of int64, in assembly where there is assembly for
// this architecture and in Go where there is not.
package sum

// Sink keeps the result reachable so the allocation test measures the call
// rather than dead code elimination.
var Sink int64

// SumInt64 returns the sum of xs.
//
// The implementation lives in sum_amd64.s. This file declares it with no body,
// which is how a Go declaration says "the definition is in assembly".
//
//go:noescape
func SumInt64(xs []int64) int64
