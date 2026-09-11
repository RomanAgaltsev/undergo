// Package drill — C2/06 shared-append.
package drill

// defaults is a base set of tags. It is created with spare capacity.
var defaults = append(make([]string, 0, 8), "a", "b", "c")

// With returns the default tags plus one extra tag.
func With(extra string) []string {
	return append(defaults, extra)
}
