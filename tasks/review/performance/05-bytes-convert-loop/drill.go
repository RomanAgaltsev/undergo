// Package drill — C7/05 bytes-convert-loop.
package drill

import "strings"

// Filter returns the lines that are a prefix of data.
func Filter(data []byte, lines []string) []string {
	var out []string
	for _, line := range lines {
		if strings.HasPrefix(string(data), line) {
			out = append(out, line)
		}
	}
	return out
}
