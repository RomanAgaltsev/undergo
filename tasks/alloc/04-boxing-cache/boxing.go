// Package boxing puts integers into an interface and asks what that costs.
//
// You do not implement anything here. Work out the answers, write them into
// prediction.yaml, then verify.
package boxing

// Sink keeps boxed values reachable so nothing is optimised away.
var Sink any

// Box puts an int into an any. That is the whole function: whatever it costs is
// the cost of the conversion itself.
func Box(n int) any { return n }
