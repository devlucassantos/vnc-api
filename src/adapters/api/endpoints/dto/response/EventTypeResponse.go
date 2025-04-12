package response

import (
	"github.com/devlucassantos/vnc-domains/src/domains/eventtype"
	"github.com/google/uuid"
)

type EventType struct {
	Id          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	Color       string    `json:"color"`
}

func NewEventType(eventType eventtype.EventType) *EventType {
	return &EventType{
		Id:          eventType.Id(),
		Description: eventType.Description(),
		Color:       eventType.Color(),
	}
}
