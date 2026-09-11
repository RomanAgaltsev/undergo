// Package drill — C4/02 leaked-cancel.
package drill

import (
	"context"
	"database/sql"
	"time"
)

// Query runs q against db with a 2-second timeout and returns the first column
// of every row.
func Query(parent context.Context, db *sql.DB, q string) ([]string, error) {
	ctx, cancel := context.WithTimeout(parent, 2*time.Second)

	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []string
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	cancel()
	return out, rows.Err()
}
