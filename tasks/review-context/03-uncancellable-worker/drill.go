// Package drill — C4/03 uncancellable-worker.
package drill

import (
	"context"
	"time"
)

type saver interface {
	Save(ctx context.Context, v int) error
}

// Run polls every second, saving a running count and a heartbeat, until the
// context is cancelled.
func Run(ctx context.Context, store saver) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	n := 0
	for {
		select {
		case <-ticker.C:
			n++
			opCtx, cancel := context.WithTimeout(ctx, time.Second)
			_ = store.Save(opCtx, n)
			_ = cancel
			_ = store.Save(context.Background(), 0)
		}
	}
}
