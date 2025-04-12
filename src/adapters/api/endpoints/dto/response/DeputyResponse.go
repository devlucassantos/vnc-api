package response

import (
	"github.com/devlucassantos/vnc-domains/src/domains/deputy"
	"github.com/google/uuid"
)

type Deputy struct {
	Id                    uuid.UUID `json:"id"`
	Name                  string    `json:"name"`
	ElectoralName         string    `json:"electoral_name"`
	ImageUrl              string    `json:"image_url"`
	ImageDescription      string    `json:"image_description"`
	Party                 *Party    `json:"party"`
	FederatedUnit         string    `json:"federated_unit"`
	PreviousParty         *Party    `json:"previous_party,omitempty"`
	PreviousFederatedUnit string    `json:"previous_federated_unit,omitempty"`
}

func NewDeputy(deputy deputy.Deputy) *Deputy {
	var previousParty *Party
	var previousFederatedUnit string
	deputyPreviousParty := deputy.PreviousParty()
	if !deputyPreviousParty.IsZero() {
		previousParty = NewParty(deputyPreviousParty)
		previousFederatedUnit = deputy.PreviousFederatedUnit()
	}

	deputyData := &Deputy{
		Id:                    deputy.Id(),
		Name:                  deputy.Name(),
		ElectoralName:         deputy.ElectoralName(),
		ImageUrl:              deputy.ImageUrl(),
		ImageDescription:      deputy.ImageDescription(),
		Party:                 NewParty(deputy.Party()),
		FederatedUnit:         deputy.FederatedUnit(),
		PreviousParty:         previousParty,
		PreviousFederatedUnit: previousFederatedUnit,
	}

	return deputyData
}
