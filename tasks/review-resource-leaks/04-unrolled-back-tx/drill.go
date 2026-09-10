// Package drill — C5/04 unrolled-back-tx.
package drill

import "database/sql"

// Transfer moves amount between two accounts in a single transaction.
func Transfer(db *sql.DB, from, to, amount int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec("UPDATE accounts SET bal = bal - ? WHERE id = ?", amount, from); err != nil {
		return err
	}
	if _, err := tx.Exec("UPDATE accounts SET bal = bal + ? WHERE id = ?", amount, to); err != nil {
		return err
	}
	return tx.Commit()
}
