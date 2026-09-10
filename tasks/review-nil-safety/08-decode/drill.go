// Package drill — C2/08 decode.
package drill

import "strconv"

// User is a decoded user.
type User struct {
	Name string
}

// Decode turns a loosely-typed payload into a string map.
func Decode(payload map[string]any) map[string]string {
	var out map[string]string

	id := payload["id"].(int)
	out["id"] = strconv.Itoa(id)

	user := payload["user"].(*User)
	out["name"] = user.Name

	return out
}

// safeString extracts a string value if present and of the right type.
func safeString(payload map[string]any, key string) string {
	if v, ok := payload[key].(string); ok {
		return v
	}
	return ""
}
