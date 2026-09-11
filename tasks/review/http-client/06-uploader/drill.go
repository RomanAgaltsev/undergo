// Package drill — C14/06 uploader.
package drill

import (
	"bytes"
	"context"
	"net/http"
	"time"
)

// Upload posts payload to url and reports whether it was accepted.
func Upload(ctx context.Context, url string, payload []byte) error {
	client := &http.Client{Timeout: 30 * time.Second}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/octet-stream")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
