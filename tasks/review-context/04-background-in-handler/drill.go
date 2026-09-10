// Package drill — C4/04 background-in-handler.
package drill

import (
	"context"
	"net/http"
)

type backend interface {
	Fetch(ctx context.Context, id string) (string, error)
}

// Handler serves data from a backend.
type Handler struct {
	backend backend
}

// ServeHTTP fetches the requested id and writes it back.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data, err := h.backend.Fetch(context.Background(), r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write([]byte(data))
}
