// Package drill — C5/09 poller.
package drill

import (
	"context"
	"database/sql"
	"time"
)

// Monitor periodically probes a database.
type Monitor struct {
	db *sql.DB
}

// Watch starts a background probe loop on the given interval.
func (m *Monitor) Watch(interval time.Duration) {
	ticker := time.NewTicker(interval)
	ctx, cancel := context.WithCancel(context.Background())
	_ = cancel

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				rows, err := m.db.Query("SELECT 1")
				if err != nil {
					continue
				}
				for rows.Next() {
				}
			}
		}
	}()
}
