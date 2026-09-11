// Package drill — C15/02 boxed-nil-pointer.
package drill

// User is a stored user.
type User struct {
	ID   int
	Name string
}

// Lookup returns the user with the given id, or an untyped nil if not found.
func Lookup(users map[int]*User, id int) any {
	var u *User
	if found, ok := users[id]; ok {
		u = found
	}
	return u
}

// NameOf returns the user's name, or "unknown" if Lookup found nothing.
func NameOf(users map[int]*User, id int) string {
	v := Lookup(users, id)
	if v == nil {
		return "unknown"
	}
	return v.(*User).Name
}
