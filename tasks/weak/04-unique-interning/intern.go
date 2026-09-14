// Package intern deduplicates repeated strings with unique.Handle.
//
// The workload is the one that makes interning pay: a stream of records whose
// label strings come from a small fixed set and repeat millions of times. Stored
// as strings, every record keeps its own copy alive. Stored as handles, they all
// point at one.
package intern

import "unique"

// Tag is a handle to an interned string.
//
// It is an alias, not a defined type, so callers can compare two Tags with ==
// and get the comparison unique.Handle promises.
type Tag = unique.Handle[string]

// Intern returns the canonical handle for s.
//
// Two calls with equal strings must return handles that compare equal, whether
// or not the two strings share a backing array.
func Intern(s string) Tag {
	panic("undergo: implement Intern")
}

// Value returns the string a handle refers to.
func Value(t Tag) string {
	panic("undergo: implement Value")
}

// InternAll returns a handle per input string, in order.
func InternAll(ss []string) []Tag {
	panic("undergo: implement InternAll")
}
