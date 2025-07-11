package response

import (
	"github.com/devlucassantos/vnc-domains/src/domains/legislativebody"
	"github.com/google/uuid"
)

type LegislativeBody struct {
	Id      uuid.UUID            `json:"id"`
	Name    string               `json:"name"`
	Acronym string               `json:"acronym"`
	Type    *LegislativeBodyType `json:"type"`
}

func NewLegislativeBody(legislativeBody legislativebody.LegislativeBody) *LegislativeBody {
	return &LegislativeBody{
		Id:      legislativeBody.Id(),
		Name:    legislativeBody.Name(),
		Acronym: legislativeBody.Acronym(),
		Type:    NewLegislativeBodyType(legislativeBody.Type()),
	}
}
