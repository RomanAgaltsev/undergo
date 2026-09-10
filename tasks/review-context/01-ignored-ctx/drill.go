// Package drill — C4/01 ignored-ctx.
package drill

import (
	"context"
	"net/http"
)

// Fetch performs a GET for url. It accepts a context for cancellation/timeout.
func Fetch(ctx context.Context, url string) (*http.Response, error) {
	return http.Get(url)
}

// FetchWithCtx performs a GET for url, honoring the context.
func FetchWithCtx(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req)
}
