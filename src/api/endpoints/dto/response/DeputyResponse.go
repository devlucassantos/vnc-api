package response

import (
	"github.com/devlucassantos/vnc-domains/src/domains/deputy"
	"github.com/google/uuid"
	"time"
)

type Deputy struct {
	Id                    uuid.UUID `json:"id,omitempty"`
	Code                  int       `json:"code,omitempty"`
	Cpf                   string    `json:"cpf,omitempty"`
	Name                  string    `json:"name,omitempty"`
	ElectoralName         string    `json:"electoral_name,omitempty"`
	ImageUrl              string    `json:"image_url,omitempty"`
	CreatedAt             time.Time `json:"created_at,omitempty"`
	UpdatedAt             time.Time `json:"updated_at,omitempty"`
	Party                 *Party    `json:"party,omitempty"`
	PartyInTheProposition *Party    `json:"party_in_the_proposal,omitempty"`
}

func NewDeputy(deputy deputy.Deputy) *Deputy {
	deputyResponse := &Deputy{
		Id:            deputy.Id(),
		Code:          deputy.Code(),
		Cpf:           deputy.Cpf(),
		Name:          deputy.Name(),
		ElectoralName: deputy.ElectoralName(),
		ImageUrl:      deputy.ImageUrl(),
		CreatedAt:     deputy.CreatedAt(),
		UpdatedAt:     deputy.CreatedAt(),
		Party:         NewParty(deputy.Party()),
	}

	partyInTheProposition := deputy.PartyInTheProposition()
	if !partyInTheProposition.IsZero() {
		deputyResponse.PartyInTheProposition = NewParty(partyInTheProposition)
	}

	return deputyResponse
}
