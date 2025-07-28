package main

type EventKind string

type Event struct {
	Kind    EventKind
	Payload any
}
