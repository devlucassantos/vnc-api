package response

import (
	"github.com/devlucassantos/vnc-domains/src/domains/eventsituation"
	"github.com/google/uuid"
)

type EventSituation struct {
	Id          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	Color       string    `json:"color"`
}

func NewEventSituation(eventSituation eventsituation.EventSituation) *EventSituation {
	return &EventSituation{
		Id:          eventSituation.Id(),
		Description: eventSituation.Description(),
		Color:       eventSituation.Color(),
	}
}
