// Package drill — C11/07 any-instantiation.
package drill

// Item is a catalog item.
type Item struct {
	ID   int
	Name string
}

// Transform maps each element of s through f.
func Transform[T, U any](s []T, f func(T) U) []U {
	out := make([]U, 0, len(s))
	for _, v := range s {
		out = append(out, f(v))
	}
	return out
}

// CollectIDs returns the ids of all items.
func CollectIDs(items []Item) []any {
	return Transform[Item, any](items, func(it Item) any {
		return it.ID
	})
}
