// Package drill — C15/05 plugin-registry.
package drill

// Handler processes a named event.
type Handler interface {
	Handle(event string) string
}

// jsonHandler is the built-in handler.
type jsonHandler struct{}

func (h *jsonHandler) Handle(event string) string { return "json:" + event }

// Registry maps plugin names to handlers.
type Registry struct {
	plugins map[string]Handler
}

// Register adds a handler under name.
func (r *Registry) Register(name string, h Handler) {
	r.plugins[name] = h
}

// Get returns the handler registered under name, or nil if none.
func (r *Registry) Get(name string) Handler {
	var h *jsonHandler
	if p, ok := r.plugins[name]; ok {
		return p
	}
	return h
}

// Default returns the built-in handler.
func (r *Registry) Default() Handler {
	var h *jsonHandler
	return h
}
