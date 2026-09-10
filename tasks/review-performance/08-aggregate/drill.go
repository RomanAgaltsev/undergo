// Package drill — C7/08 aggregate.
package drill

import "strings"

// Aggregate prefixes each chunk with the header and joins them.
func Aggregate(header []byte, chunks [][]byte) string {
	var parts []string
	for _, c := range chunks {
		defer trace()
		line := string(header) + ":" + string(c)
		parts = append(parts, line)
	}
	return strings.Join(parts, "\n")
}

func trace() {}
