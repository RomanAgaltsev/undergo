// Package drill — C12/06 config-loader.
package drill

import "encoding/json"

// Config is loaded from a JSON config document.
type Config struct {
	Enabled   bool   `json:"enabled,omitempty"`
	secret    string `json:"secret"`
	UpdatedAt int64  `json:"-"`
}

// LoadConfig decodes a Config and resolves its update timestamp from the raw payload.
func LoadConfig(data []byte) (Config, error) {
	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		return Config{}, err
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return Config{}, err
	}
	c.UpdatedAt = int64(raw["updated_at"].(float64))

	return c, nil
}

// Secret returns the loaded secret.
func (c Config) Secret() string { return c.secret }
