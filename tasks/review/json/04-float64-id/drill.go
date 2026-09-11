// Package drill — C12/04 float64-id.
package drill

import "encoding/json"

// ParseOrderID decodes a JSON payload and returns the order id.
func ParseOrderID(data []byte) (int64, error) {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return 0, err
	}
	id := int64(m["id"].(float64))
	return id, nil
}
