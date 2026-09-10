// Package drill — C15/08 slice-of-handlers.
package drill

// Handler processes an event.
type Handler interface {
	Handle(event string)
}

// fnHandler is a handler backed by a function.
type fnHandler struct {
	fn func(string)
}

func (h *fnHandler) Handle(event string) { h.fn(event) }

// newHandler returns the handler for name, or a nil *fnHandler for unknown names.
func newHandler(name string, table map[string]func(string)) *fnHandler {
	fn, ok := table[name]
	if !ok {
		return nil
	}
	return &fnHandler{fn: fn}
}

// Dispatch builds handlers for the given names and runs each against event.
func Dispatch(names []string, table map[string]func(string), event string) {
	handlers := make([]Handler, 0, len(names))
	for _, name := range names {
		handlers = append(handlers, newHandler(name, table))
	}

	for _, h := range handlers {
		if h != nil {
			h.Handle(event)
		}
	}
}
