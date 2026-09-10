// Package padding holds the types this task asks you to measure.
//
// You do not implement anything here. Read the types, work out the answers,
// and write them into prediction.yaml. Then verify.
package padding

// Header is laid out in declaration order, with padding wherever a field's
// alignment demands it.
type Header struct {
	Flag    bool
	Offset  int64
	Kind    uint16
	Payload []byte
}

// Reordered holds exactly the same four fields.
type Reordered struct {
	Payload []byte
	Offset  int64
	Kind    uint16
	Flag    bool
}

// Trailing ends with a zero-size field. That is not free.
type Trailing struct {
	Count int64
	Done  struct{}
}
