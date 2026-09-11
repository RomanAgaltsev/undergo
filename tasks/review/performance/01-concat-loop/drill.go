// Package drill — C7/01 concat-loop.
package drill

import "strings"

// Join concatenates parts into a single string.
func Join(parts []string) string {
	s := ""
	for _, p := range parts {
		s += p
	}
	return s
}

// JoinBuilder concatenates parts using a strings.Builder.
func JoinBuilder(parts []string) string {
	var b strings.Builder
	for _, p := range parts {
		b.WriteString(p)
	}
	return b.String()
}
