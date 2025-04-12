package services

import (
	"github.com/devlucassantos/vnc-domains/src/domains/event"
	"github.com/google/uuid"
	"vnc-api/core/interfaces/postgres"
)

type Event struct {
	repository postgres.Event
}

func NewEventService(repository postgres.Event) *Event {
	return &Event{
		repository: repository,
	}
}

func (instance Event) GetEventByArticleId(articleId uuid.UUID, userId uuid.UUID) (*event.Event, error) {
	return instance.repository.GetEventByArticleId(articleId, userId)
}
