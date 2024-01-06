package response

import (
	"github.com/google/uuid"
	"time"
	"vnc-read-api/core/domains/deputy"
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
	return &Deputy{
		Id:                    deputy.Id(),
		Code:                  deputy.Code(),
		Cpf:                   deputy.Cpf(),
		Name:                  deputy.Name(),
		ElectoralName:         deputy.ElectoralName(),
		ImageUrl:              deputy.ImageUrl(),
		CreatedAt:             deputy.CreatedAt(),
		UpdatedAt:             deputy.CreatedAt(),
		Party:                 NewParty(deputy.Party()),
		PartyInTheProposition: NewParty(deputy.PartyInTheProposition()),
	}
}
