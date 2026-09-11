// Package snippets holds five small functions. Your job is to predict how many
// heap allocations each performs per call.
//
// You do not implement anything here. Work out the answers, write them into
// prediction.yaml, then verify.
package snippets

import (
	"fmt"
	"strconv"
	"strings"
)

// The sinks keep results reachable so nothing is optimised away. They are typed
// rather than any, because assigning to an any would box the result and add an
// allocation that belongs to the measurement rather than to the function.
var (
	SinkString string
	SinkInt    int
)

// ConcatTwo joins two strings with +.
func ConcatTwo(a, b string) string { return a + b }

// SprintfInt formats an int with fmt.
func SprintfInt(n int) string { return fmt.Sprintf("%d", n) }

// ItoaInt formats the same int with strconv.
func ItoaInt(n int) string { return strconv.Itoa(n) }

// BuilderJoin joins a pair of strings with a strings.Builder.
func BuilderJoin(a, b string) string {
	var sb strings.Builder
	sb.WriteString(a)
	sb.WriteString(b)
	return sb.String()
}

// MakeLocal allocates a fixed-size slice that never leaves the function.
func MakeLocal() int {
	buf := make([]byte, 64)
	return len(buf)
}
