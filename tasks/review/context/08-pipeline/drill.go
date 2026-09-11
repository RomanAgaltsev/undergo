// Package drill — C4/08 pipeline.
package drill

import (
	"context"
	"time"
)

// Pipeline enriches the context, fetches data, then notifies downstream.
func Pipeline(
	ctx context.Context,
	id string,
	fetch func(context.Context, string) (string, error),
	notify func(context.Context, string) error,
) error {
	ctx = context.WithValue(ctx, "requestID", id)

	data, err := fetch(context.Background(), id)
	if err != nil {
		return err
	}

	nctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return notify(nctx, data)
}
