// Package drill — C5/03 ticker-rows.
package drill

import (
	"context"
	"database/sql"
	"time"
)

// Poller periodically counts rows in a table.
type Poller struct {
	db *sql.DB
}

// Run starts polling every interval in the background and returns a stop function.
func (p *Poller) Run(interval time.Duration) func() {
	ticker := time.NewTicker(interval)
	ctx, cancel := context.WithCancel(context.Background())
	_ = cancel

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				p.countOnce()
			}
		}
	}()

	return func() {
		// intended to stop the poller
	}
}

func (p *Poller) countOnce() {
	rows, err := p.db.Query("SELECT id FROM items")
	if err != nil {
		return
	}
	for rows.Next() {
		var id int
		_ = rows.Scan(&id)
	}
}

// queryStmt prepares and closes a statement correctly.
func (p *Poller) queryStmt() error {
	stmt, err := p.db.Prepare("SELECT 1")
	if err != nil {
		return err
	}
	defer stmt.Close()
	return nil
}
