// Package join asks you to build a string join that allocates nothing when the
// caller supplies enough capacity.
package join

import "strings"

// joinBaseline is the naive version, kept as the benchmark's reference point.
// Do not change it.
func joinBaseline(parts []string, sep string) string {
	out := ""
	for i, p := range parts {
		if i > 0 {
			out += sep
		}
		out += p
	}
	return out
}

// JoinInto appends parts, separated by sep, to dst and returns the result.
//
// When dst has enough capacity for the whole result, JoinInto must not allocate
// at all. The returned slice must contain exactly what strings.Join(parts, sep)
// would produce.
func JoinInto(dst []byte, parts []string, sep string) []byte {
	// Replace this. It is correct and it allocates.
	return append(dst, strings.Join(parts, sep)...)
}
