package response

import (
	"github.com/devlucassantos/vnc-domains/src/domains/legislativebodytype"
	"github.com/google/uuid"
)

type LegislativeBodyType struct {
	Id          uuid.UUID `json:"id"`
	Description string    `json:"description"`
}

func NewLegislativeBodyType(legislativeBodyType legislativebodytype.LegislativeBodyType) *LegislativeBodyType {
	return &LegislativeBodyType{
		Id:          legislativeBodyType.Id(),
		Description: legislativeBodyType.Description(),
	}
}
