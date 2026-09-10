// Package drill — C3/02 lost-wrap.
package drill

import (
	"database/sql"
	"errors"
	"fmt"
)

// ErrNotFound is returned when a user does not exist.
var ErrNotFound = errors.New("user not found")

// Store loads users from a database.
type Store struct{ db *sql.DB }

// LoadUser returns the user's name, or ErrNotFound if there is no such user.
func (s *Store) LoadUser(id int) (string, error) {
	var name string
	err := s.db.QueryRow("SELECT name FROM users WHERE id = ?", id).Scan(&name)
	if errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("load user %d: %v", id, ErrNotFound)
	}
	if err != nil {
		return "", fmt.Errorf("load user %d: %w", id, err)
	}
	return name, nil
}

// StatusFor maps a LoadUser error to an HTTP status code.
func StatusFor(err error) int {
	switch {
	case err == nil:
		return 200
	case errors.Is(err, ErrNotFound):
		return 404
	default:
		return 500
	}
}
