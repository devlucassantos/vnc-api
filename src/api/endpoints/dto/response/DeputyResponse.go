package response

import (
	"github.com/devlucassantos/vnc-domains/src/domains/deputy"
	"github.com/google/uuid"
	"time"
)

type Deputy struct {
	Id                    uuid.UUID `json:"id"`
	Name                  string    `json:"name"`
	ElectoralName         string    `json:"electoral_name"`
	ImageUrl              string    `json:"image_url"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
	Party                 *Party    `json:"party"`
	PartyInTheProposition *Party    `json:"party_in_the_proposal,omitempty"`
}

func NewDeputy(deputy deputy.Deputy) *Deputy {
	deputyData := &Deputy{
		Id:            deputy.Id(),
		Name:          deputy.Name(),
		ElectoralName: deputy.ElectoralName(),
		ImageUrl:      deputy.ImageUrl(),
		CreatedAt:     deputy.CreatedAt(),
		UpdatedAt:     deputy.CreatedAt(),
		Party:         NewParty(deputy.Party()),
	}

	deputyPartyInTheProposition := deputy.PartyInTheProposition()
	if deputyPartyInTheProposition.Id() != uuid.Nil {
		deputyData.PartyInTheProposition = NewParty(deputyPartyInTheProposition)
	}

	return deputyData
}
