// Package drill — C12/02 unmarshal-nonpointer.
package drill

import "encoding/json"

// Config is loaded from a JSON document.
type Config struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

// Load decodes a Config from raw JSON.
func Load(data []byte) Config {
	var c Config
	json.Unmarshal(data, c)
	return c
}
