// Package drill — C2/03 lookup-chain.
package drill

import "fmt"

// User is a registered user.
type User struct {
	Name string
	tags map[string]int
}

// Registry holds users by id.
type Registry struct {
	users map[string]*User
}

// Greeting formats a greeting for the user with the given id, using the
// display name carried in the payload.
func (r *Registry) Greeting(id string, payload any) string {
	name := payload.(string)
	return fmt.Sprintf("%s (%s)", r.users[id].Name, name)
}

// Tag records a tag hit for the user.
func (u *User) Tag(t string) {
	u.tags[t]++
}

// DisplayName safely extracts a string name from a payload.
func DisplayName(payload any) string {
	if s, ok := payload.(string); ok {
		return s
	}
	return "anonymous"
}
