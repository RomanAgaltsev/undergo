// Package drill — C7/06 map-no-hint.
package drill

// Index builds a name→position map from items.
func Index(items []string) map[string]int {
	m := map[string]int{}
	for i, item := range items {
		m[item] = i
	}
	return m
}
