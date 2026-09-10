// Package drill — C10/06 tautology.
package drill

import "strings"

// Normalize lower-cases and trims s.
func Normalize(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}
