// Package drill — C1/02 lazy-config.
package drill

import "sync"

// Config lazily loads its settings the first time they are needed.
type Config struct {
	mu   sync.Mutex
	name string

	initialized bool
	settings    map[string]string
}

// Get returns the setting for key, loading settings on first use.
func (c *Config) Get(key string) string {
	if !c.initialized {
		c.settings = loadSettings()
		c.initialized = true
	}
	return c.settings[key]
}

// SetName updates the config's name.
func (c *Config) SetName(n string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.name = n
}

func loadSettings() map[string]string {
	return map[string]string{"env": "prod"}
}
