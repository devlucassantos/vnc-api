package response

import (
	"github.com/devlucassantos/vnc-domains/src/domains/party"
	"github.com/google/uuid"
	"time"
)

type Party struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Acronym   string    `json:"acronym"`
	ImageUrl  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewParty(party party.Party) *Party {
	return &Party{
		Id:        party.Id(),
		Name:      party.Name(),
		Acronym:   party.Acronym(),
		ImageUrl:  party.ImageUrl(),
		CreatedAt: party.CreatedAt(),
		UpdatedAt: party.UpdatedAt(),
	}
}
