package events

type EventKind string

type Event struct {
	Kind    EventKind
	Payload any
}
