package response

import (
	"github.com/devlucassantos/vnc-domains/src/domains/externalauthor"
	"github.com/google/uuid"
)

type ExternalAuthor struct {
	Id   uuid.UUID           `json:"id"`
	Name string              `json:"name"`
	Type *ExternalAuthorType `json:"type"`
}

func NewExternalAuthor(externalAuthor externalauthor.ExternalAuthor) *ExternalAuthor {
	return &ExternalAuthor{
		Id:   externalAuthor.Id(),
		Name: externalAuthor.Name(),
		Type: NewExternalAuthorType(externalAuthor.Type()),
	}
}
