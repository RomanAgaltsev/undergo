// Package drill — C12/09 message.
package drill

import "encoding/json"

// Meta is optional per-message metadata.
type Meta struct {
	Client string `json:"client"`
}

// Message is a chat message decoded from the wire.
type Message struct {
	Body      string `json:"body"`
	SessionID string `json:",omitempty"`
	Meta      *Meta  `json:"meta"`
}

// Decode parses a message and returns its body plus the reporting client.
func Decode(data []byte) (Message, string, error) {
	var m Message
	json.Unmarshal(data, m)

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return Message{}, "", err
	}
	ts := int64(raw["timestamp"].(float64))
	_ = ts

	return m, m.Meta.Client, nil
}
