// Package drill — C1/01 request-counter.
package drill

import "net/http"

// Server counts how many times each URL path has been requested.
type Server struct {
	counts map[string]int
}

// NewServer returns a ready Server.
func NewServer() *Server {
	return &Server{counts: make(map[string]int)}
}

// ServeHTTP records a hit for the request path. To avoid making the caller wait
// on the bookkeeping, the increment is done in a background goroutine.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	go func() {
		s.counts[path]++
	}()
	w.WriteHeader(http.StatusAccepted)
}

// Counts returns the recorded hit counts.
func (s *Server) Counts() map[string]int {
	return s.counts
}
