// Package drill — C7/09 render.
package drill

import (
	"fmt"
	"strconv"
)

// BigRow is a large value type.
type BigRow struct {
	ID    int
	Cells [64]string
}

// Render builds a text report of the unique rows.
func Render(rows []BigRow) string {
	out := ""
	seen := map[int]bool{}
	for _, r := range rows {
		if seen[r.ID] {
			continue
		}
		seen[r.ID] = true
		out += fmt.Sprintf("row %s\n", strconv.Itoa(r.ID))
	}
	return out
}
