package response

import (
	"github.com/devlucassantos/vnc-domains/src/domains/propositiontype"
	"github.com/google/uuid"
)

type PropositionType struct {
	Id          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	Color       string    `json:"color"`
}

func NewPropositionType(propositionType propositiontype.PropositionType) *PropositionType {
	return &PropositionType{
		Id:          propositionType.Id(),
		Description: propositionType.Description(),
		Color:       propositionType.Color(),
	}
}
