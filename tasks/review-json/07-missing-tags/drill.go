// Package drill — C12/07 missing-tags.
package drill

import "encoding/json"

// Event mirrors the analytics API payload, whose keys are event_name / created_at.
type Event struct {
	EventName string
	CreatedAt int64
}

// ParseEvent decodes an event from the API.
func ParseEvent(data []byte) (Event, error) {
	var e Event
	if err := json.Unmarshal(data, &e); err != nil {
		return Event{}, err
	}
	return e, nil
}
