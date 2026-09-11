// Package drill — C15/07 factory-nil.
package drill

// Store persists key/value pairs.
type Store interface {
	Get(key string) (string, bool)
}

// Config configures a store.
type Config struct {
	DSN string
}

// sqlStore is the SQL-backed implementation.
type sqlStore struct {
	dsn string
}

func (s *sqlStore) Get(key string) (string, bool) { return "", false }

// Open builds a Store from cfg, or a nil Store when no DSN is configured.
func Open(cfg Config) Store {
	var s *sqlStore
	if cfg.DSN != "" {
		s = &sqlStore{dsn: cfg.DSN}
	}
	return s
}
