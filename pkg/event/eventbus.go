package event

const (
	PostVisitedEvent = "post_visited"
)

type Event struct {
	Type string
	Data any
}

type EventBus struct {
	bus chan Event
}

func NewEventBus() *EventBus {
	return &EventBus{
		bus: make(chan Event),
	}
}

func (eb *EventBus) Send(e Event) {
	eb.bus <- e
}

func (eb *EventBus) Subscribe() <-chan Event {
	return eb.bus
}
