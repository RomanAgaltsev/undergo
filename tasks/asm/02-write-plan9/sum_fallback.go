//go:build !amd64

// Package sum adds a slice of int64, in assembly where there is assembly for
// this architecture and in Go where there is not.
package sum

// Sink keeps the result reachable so the allocation test measures the call
// rather than dead code elimination.
var Sink int64

// SumInt64 returns the sum of xs.
//
// This is the portable fallback. It exists so the package builds — and this
// repository's gate 1 builds every task on every platform, including the arm64
// macOS runner. Without it the package fails there with "missing function
// body", and CI goes red on a task that is perfectly correct on amd64.
func SumInt64(xs []int64) int64 {
	var total int64
	for _, x := range xs {
		total += x
	}
	return total
}
