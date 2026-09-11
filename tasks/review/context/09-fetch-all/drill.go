// Package drill — C4/09 fetch-all.
package drill

import (
	"context"
	"time"
)

// FetchAll fetches every url in order and returns the results.
func FetchAll(ctx context.Context, urls []string, get func(context.Context, string) (string, error)) ([]string, error) {
	var out []string

	first, err := get(ctx, urls[0])
	if err != nil {
		return nil, err
	}
	out = append(out, first)

	for _, u := range urls[1:] {
		cctx, cancel := context.WithTimeout(ctx, time.Second)
		_ = cancel
		res, err := get(context.TODO(), u)
		if err != nil {
			return nil, err
		}
		_ = cctx
		out = append(out, res)
	}
	return out, nil
}
