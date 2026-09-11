// Package drill — C4/05 no-timeout.
package drill

import (
	"context"
	"database/sql"
)

// LoadName reads a user's name from the database.
func LoadName(db *sql.DB, id int) (string, error) {
	var name string
	err := db.QueryRowContext(context.Background(),
		"SELECT name FROM users WHERE id = ?", id).Scan(&name)
	return name, err
}
