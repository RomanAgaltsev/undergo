// Package drill — C8/01 sql-injection.
package drill

import "database/sql"

// FindUser looks up a user id by name.
func FindUser(db *sql.DB, name string) *sql.Row {
	q := "SELECT id FROM users WHERE name = '" + name + "'"
	return db.QueryRow(q)
}

// FindByID looks up a user name by id.
func FindByID(db *sql.DB, id int) *sql.Row {
	return db.QueryRow("SELECT name FROM users WHERE id = ?", id)
}
