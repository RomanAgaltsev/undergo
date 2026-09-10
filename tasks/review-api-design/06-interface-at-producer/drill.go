// Package drill — C6/06 interface-at-producer.
package drill

// Store is the key-value store abstraction exposed by this package.
type Store interface {
	Get(key string) (string, bool)
	Set(key, value string)
}

type memStore struct {
	data map[string]string
}

// NewStore returns a new in-memory store.
func NewStore() Store {
	return &memStore{data: make(map[string]string)}
}

func (m *memStore) Get(key string) (string, bool) {
	v, ok := m.data[key]
	return v, ok
}

func (m *memStore) Set(key, value string) {
	m.data[key] = value
}
