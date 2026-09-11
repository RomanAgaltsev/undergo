// Package drill — C8/08 download.
package drill

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// Download fetches url and saves the body as name inside dir.
func Download(dir, name, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	path := filepath.Join(dir, name)
	return os.WriteFile(path, data, 0o644)
}
