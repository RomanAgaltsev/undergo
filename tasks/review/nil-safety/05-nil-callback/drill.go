// Package drill — C2/05 nil-callback.
package drill

// Server emits events to an optional callback.
type Server struct {
	Name    string
	OnEvent func(string)
}

// Emit delivers msg to the event callback.
func (s *Server) Emit(msg string) {
	s.OnEvent(msg)
}
