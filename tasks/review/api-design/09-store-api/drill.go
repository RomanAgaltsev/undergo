// Package drill — C6/09 store-api.
package drill

type record struct {
	id   int
	name string
}

// Store keeps records by id.
type Store struct {
	records map[int]*record
	name    string
}

// NewStore returns a named, empty store.
func NewStore(name string) *Store {
	return &Store{records: make(map[int]*record), name: name}
}

// GetName returns the store's name.
func (s *Store) GetName() string {
	return s.name
}

// Find returns the record for id.
func (s *Store) Find(id int) *record {
	return s.records[id]
}

// Update replaces a record's fields.
func (s *Store) Update(id int, name string, active bool, score int, tags []string) {
	s.records[id] = &record{id: id, name: name}
}
