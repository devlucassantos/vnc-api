package postgres

import (
	"github.com/devlucassantos/vnc-domains/src/domains/event"
	"github.com/google/uuid"
)

type Event interface {
	GetEventByArticleId(articleId uuid.UUID, userId uuid.UUID) (*event.Event, error)
}
