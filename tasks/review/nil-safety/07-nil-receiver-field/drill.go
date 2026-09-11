// Package drill — C2/07 nil-receiver-field.
package drill

// List is an integer list.
type List struct {
	items []int
}

// Len returns the number of items. It is documented as safe to call on a nil *List.
func (l *List) Len() int {
	return len(l.items)
}
