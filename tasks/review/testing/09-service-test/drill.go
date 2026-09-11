// Package drill — C10/09 service-test.
package drill

import "errors"

// Record is a stored record.
type Record struct {
	Name string
}

// Service assigns incrementing ids to saved records.
type Service struct {
	nextID int
}

// Save stores rec and returns its new id, erroring on a missing name.
func (s *Service) Save(rec Record) (int, error) {
	if rec.Name == "" {
		return 0, errors.New("name required")
	}
	s.nextID++
	return s.nextID, nil
}
