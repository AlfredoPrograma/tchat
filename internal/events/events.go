package events

import (
	"encoding/json"

	"github.com/alfredoprograma/tchat/internal/log"
)

type EventKind string

const (
	REGISTER_USER_EVENT = "REGISTER_USER_EVENT"
	SEND_MESSAGE_EVENT  = "SEND_MESSAGE_EVENT"
)

type Event struct {
	Kind    EventKind
	Payload any
}

type RegisterUserPayload struct {
	username string
}

func NewRegisterUserEvent(username string) Event {
	return Event{
		Kind: REGISTER_USER_EVENT,
		Payload: RegisterUserPayload{
			username: username,
		},
	}
}

type SendMessagePayload struct {
	content string
}

func NewSendMessageEvent(content string) Event {
	return Event{
		Kind: SEND_MESSAGE_EVENT,
		Payload: SendMessagePayload{
			content: content,
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
	}

	return event
}
