package response

import (
	"github.com/devlucassantos/vnc-domains/src/domains/external"
	"github.com/google/uuid"
	"time"
)

type ExternalAuthor struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewExternalAuthor(externalAuthor external.ExternalAuthor) *ExternalAuthor {
	return &ExternalAuthor{
		Id:        externalAuthor.Id(),
		Name:      externalAuthor.Name(),
		Type:      externalAuthor.Type(),
		CreatedAt: externalAuthor.CreatedAt(),
		UpdatedAt: externalAuthor.UpdatedAt(),
	}
}
