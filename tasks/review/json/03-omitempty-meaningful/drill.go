// Package drill — C12/03 omitempty-meaningful.
package drill

import "encoding/json"

// Prefs are a user's saved preferences, sent to the client as JSON.
type Prefs struct {
	Notifications bool   `json:"notifications,omitempty"`
	MaxItems      int    `json:"max_items,omitempty"`
	Note          string `json:"note,omitempty"`
}

// Encode renders prefs as JSON.
func Encode(p Prefs) ([]byte, error) {
	return json.Marshal(p)
}
