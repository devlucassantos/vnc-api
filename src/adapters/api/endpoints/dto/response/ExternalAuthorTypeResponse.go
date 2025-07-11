package response

import (
	"github.com/devlucassantos/vnc-domains/src/domains/externalauthortype"
	"github.com/google/uuid"
)

type ExternalAuthorType struct {
	Id          uuid.UUID `json:"id"`
	Description string    `json:"description"`
}

func NewExternalAuthorType(externalAuthorType externalauthortype.ExternalAuthorType) *ExternalAuthorType {
	return &ExternalAuthorType{
		Id:          externalAuthorType.Id(),
		Description: externalAuthorType.Description(),
	}
}
