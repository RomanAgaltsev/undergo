// Package drill — C15/09 event-bus.
package drill

// Event is a published event.
type Event struct {
	Name string
}

// Handler processes an event.
type Handler interface {
	Handle(e Event)
}

// Unsubscriber cancels a subscription.
type Unsubscriber interface {
	Unsubscribe()
}

// subscription implements Unsubscriber.
type subscription struct {
	bus *Bus
	id  int
}

func (s *subscription) Unsubscribe() { delete(s.bus.handlers, s.id) }

// fnHandler adapts a func to Handler.
type fnHandler struct {
	fn func(Event)
}

func (h *fnHandler) Handle(e Event) { h.fn(e) }

// Bus is a minimal in-process event bus.
type Bus struct {
	handlers map[int]Handler
	nextID   int
	last     *Event
}

// NewBus builds an empty bus.
func NewBus() *Bus {
	return &Bus{handlers: make(map[int]Handler)}
}

// Subscribe registers fn and returns a handle to cancel it.
func (b *Bus) Subscribe(fn func(Event)) Unsubscriber {
	var sub *subscription
	if fn != nil {
		b.nextID++
		b.handlers[b.nextID] = &fnHandler{fn: fn}
		sub = &subscription{bus: b, id: b.nextID}
	}
	return sub
}

// find returns the fnHandler registered under id, or nil.
func (b *Bus) find(id int) Handler {
	var h *fnHandler
	if got, ok := b.handlers[id]; ok {
		return got
	}
	return h
}

// Has reports whether a handler is registered under id.
func (b *Bus) Has(id int) bool {
	return b.find(id) != nil
}

// Publish delivers e to every subscribed handler.
func (b *Bus) Publish(e Event) {
	b.last = &e
	for _, h := range b.handlers {
		h.Handle(e)
	}
}

// Last returns the most recently published event, or nil if none.
func (b *Bus) Last() any {
	return b.last
}
