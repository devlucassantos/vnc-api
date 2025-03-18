package response

import (
	"github.com/devlucassantos/vnc-domains/src/domains/agendaitemregime"
	"github.com/google/uuid"
)

type AgendaItemRegime struct {
	Id          uuid.UUID `json:"id"`
	Description string    `json:"description"`
}

func NewAgendaItemRegime(agendaItemRegime agendaitemregime.AgendaItemRegime) *AgendaItemRegime {
	return &AgendaItemRegime{
		Id:          agendaItemRegime.Id(),
		Description: agendaItemRegime.Description(),
	}
}
