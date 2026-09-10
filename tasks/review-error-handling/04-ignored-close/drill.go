// Package drill — C3/04 ignored-close.
package drill

import (
	"encoding/json"
	"os"
)

// WriteConfig writes cfg to path as JSON.
func WriteConfig(path string, cfg any) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(cfg)
}
