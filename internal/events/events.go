package events

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/alfredoprograma/tchat/internal/log"
)

const BUFFER_SIZE = 512

type EventKind string

const (
	REGISTER_USER_EVENT = "REGISTER_USER_EVENT"
	SEND_MESSAGE_EVENT  = "SEND_MESSAGE_EVENT"
)

type EventMetadata struct {
	Conn *net.TCPConn
}

type Event struct {
	Meta    EventMetadata
	Kind    EventKind
	Payload any
}

type RegisterUserPayload struct {
	Username string
}

func NewRegisterUserEvent(username string) Event {
	return Event{
		Kind: REGISTER_USER_EVENT,
		Payload: RegisterUserPayload{
			Username: username,
		},
	}
}

type SendMessagePayload struct {
	Content string
}

func NewSendMessageEvent(content string) Event {
	return Event{
		Kind: SEND_MESSAGE_EVENT,
		Payload: SendMessagePayload{
			Content: content,
		},
	}
}

func Serialize(event Event) []byte {
	raw, err := json.Marshal(event)

	if err != nil {
		log.Log(log.LOG_LEVEL_ERROR, "cannot serialize event")
		return []byte{}
	}

	return raw
}

func Deserialize(raw []byte) Event {
	var event Event

	if err := json.Unmarshal(raw, &event); err != nil {
		log.Log(log.LOG_LEVEL_ERROR, "cannot deserialize event")
		return Event{}
	}

	deserializePayload(&event)
	return event
}

// TODO: deserialize event, then serialize payload and deserialize it again to narrow its type
// is kinda inefficient.
//
// Here are some solutions I think:
//
// - Use reflection to build structs
//
// - Use a custom json unmarshaler to directly unmashal complete event and its payload
func deserializePayload(event *Event) {
	rawPayload, err := json.Marshal(event.Payload)

	if err != nil {
		log.Log(log.LOG_LEVEL_ERROR, fmt.Sprintf("invalid payload for associated event %s", event.Kind))
		return
	}

	switch event.Kind {
	case REGISTER_USER_EVENT:
		var payload RegisterUserPayload
		if err = json.Unmarshal(rawPayload, &payload); err != nil {
			log.Log(log.LOG_LEVEL_ERROR, "cannot deserialize event payload")
			return
		}

		event.Payload = payload
	case SEND_MESSAGE_EVENT:
		var payload SendMessagePayload
		if err = json.Unmarshal(rawPayload, &payload); err != nil {
			log.Log(log.LOG_LEVEL_ERROR, "cannot deserialize event payload")
			return
		}

		event.Payload = payload
	}
}
