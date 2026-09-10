// Package drill — C7/02 n-plus-one.
package drill

import "database/sql"

// LoadTotals returns the summed order amount for each user id.
func LoadTotals(db *sql.DB, userIDs []int) (map[int]int, error) {
	totals := make(map[int]int, len(userIDs))
	for _, id := range userIDs {
		var t int
		if err := db.QueryRow("SELECT SUM(amount) FROM orders WHERE user_id = ?", id).Scan(&t); err != nil {
			return nil, err
		}
		totals[id] = t
	}
	return totals, nil
}
