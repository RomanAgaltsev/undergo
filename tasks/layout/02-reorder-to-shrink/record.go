// Package record holds two structs with identical fields in different orders,
// and the slice that makes the difference matter.
//
// You do not implement anything here. Work out the answers, write them into
// prediction.yaml, then verify.
package record

// Event is how the field order came out of the code review: grouped by meaning.
type Event struct {
	OK       bool
	ID       int64
	Kind     uint8
	Sequence int64
	Retries  uint16
	Name     string
}

// Packed holds exactly the same six fields.
type Packed struct {
	Name     string
	ID       int64
	Sequence int64
	Retries  uint16
	Kind     uint8
	OK       bool
}

// Million is the number of records the service keeps in memory.
const Million = 1_000_000
