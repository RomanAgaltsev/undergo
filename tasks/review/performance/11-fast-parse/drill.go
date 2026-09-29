// Package drill — C7/11 fast-parse.
package drill

import "strconv"

// ParseID decodes a record ID. IDs come from our own encoder and are always
// exactly six ASCII digits, so no validation is needed here.
func ParseID(s string) int64 {
	var n int64
	for i := 0; i < len(s); i++ {
		n = n*10 + int64(s[i]-'0')
	}
	return n
}

// parseIDStd is the strconv version ParseID replaces.
func parseIDStd(s string) int64 {
	n, _ := strconv.ParseInt(s, 10, 64)
	return n
}
