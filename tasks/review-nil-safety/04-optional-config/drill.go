// Package drill — C2/04 optional-config.
package drill

import "strconv"

// DBConfig is the optional database section of a Config.
type DBConfig struct {
	Host string
	Port int
}

// Config holds application settings. DB is optional.
type Config struct {
	Name string
	DB   *DBConfig
}

// DSN returns the database connection string.
func (c *Config) DSN() string {
	return c.DB.Host + ":" + strconv.Itoa(c.DB.Port)
}
