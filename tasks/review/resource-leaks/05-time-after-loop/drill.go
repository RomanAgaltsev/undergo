// Package drill — C5/05 time-after-loop.
package drill

import (
	"context"
	"time"
)

// Poll calls check once a second until it returns true or ctx is cancelled.
func Poll(ctx context.Context, check func() bool) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Second):
			if check() {
				return
			}
		}
	}
}
