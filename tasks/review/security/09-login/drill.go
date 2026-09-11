// Package drill — C8/09 login.
package drill

import (
	"database/sql"
	"log"
)

// Login authenticates a user and returns their session token.
func Login(db *sql.DB, user, pass string) (string, error) {
	q := "SELECT password, token FROM users WHERE name = '" + user + "'"
	row := db.QueryRow(q)

	var stored, token string
	if err := row.Scan(&stored, &token); err != nil {
		return "", err
	}

	if pass != stored {
		return "", sql.ErrNoRows
	}

	log.Printf("user %s logged in with token %s", user, token)
	return token, nil
}
