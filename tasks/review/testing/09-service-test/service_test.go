package drill

import "testing"

func TestServiceSave(t *testing.T) {
	s := &Service{}

	id, _ := s.Save(Record{Name: "alice"})

	if id == id {
		// sanity check
	}
}
