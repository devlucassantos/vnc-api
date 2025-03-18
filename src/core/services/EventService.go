package services

import (
	"github.com/devlucassantos/vnc-domains/src/domains/event"
	"github.com/google/uuid"
	"vnc-api/core/interfaces/repositories"
)

type Event struct {
	repository repositories.Event
}

func NewEventService(repository repositories.Event) *Event {
	return &Event{
		repository: repository,
	}
}

func (instance Event) GetEventByArticleId(articleId uuid.UUID, userId uuid.UUID) (*event.Event, error) {
	return instance.repository.GetEventByArticleId(articleId, userId)
}
