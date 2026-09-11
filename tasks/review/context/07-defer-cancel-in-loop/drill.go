// Package drill — C4/07 defer-cancel-in-loop.
package drill

import (
	"context"
	"time"
)

// ProcessAll runs fn for each id, each under a one-second timeout.
func ProcessAll(ids []int, fn func(context.Context, int) error) error {
	for _, id := range ids {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		if err := fn(ctx, id); err != nil {
			return err
		}
	}
	return nil
}
