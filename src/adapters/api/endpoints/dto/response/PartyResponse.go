package response

import (
	"github.com/devlucassantos/vnc-domains/src/domains/party"
	"github.com/google/uuid"
)

type Party struct {
	Id               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	Acronym          string    `json:"acronym"`
	ImageUrl         string    `json:"image_url"`
	ImageDescription string    `json:"image_description"`
}

func NewParty(party party.Party) *Party {
	return &Party{
		Id:               party.Id(),
		Name:             party.Name(),
		Acronym:          party.Acronym(),
		ImageUrl:         party.ImageUrl(),
		ImageDescription: party.ImageDescription(),
	}
}
