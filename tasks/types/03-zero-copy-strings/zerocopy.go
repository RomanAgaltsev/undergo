// Package zerocopy converts between string and []byte without copying.
package zerocopy

// Sinks keep results reachable so the allocation test measures the conversion
// rather than dead code elimination.
var (
	Sink       []byte
	SinkString string
)

// StringToBytes returns a []byte sharing s's backing array.
//
// The result must not be written to: a string's bytes are immutable, and
// mutating them through this slice is undefined behaviour that neither the
// compiler nor the race detector will catch.
//
// It must allocate nothing, and must return a nil slice for the empty string.
func StringToBytes(s string) []byte {
	panic("undergo: implement StringToBytes")
}

// BytesToString returns a string sharing b's backing array.
//
// The result is only valid while the caller does not modify b. It must allocate
// nothing, and must return the empty string for a nil or empty slice.
func BytesToString(b []byte) string {
	panic("undergo: implement BytesToString")
}
