package response

import (
	"github.com/google/uuid"
	"time"
	"vnc-read-api/core/domains/party"
)

type Party struct {
	Id        uuid.UUID `json:"id"`
	Code      int       `json:"code"`
	Name      string    `json:"name"`
	Acronym   string    `json:"acronym"`
	ImageUrl  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewParty(party party.Party) *Party {
	return &Party{
		Id:        party.Id(),
		Code:      party.Code(),
		Name:      party.Name(),
		Acronym:   party.Acronym(),
		ImageUrl:  party.ImageUrl(),
		CreatedAt: party.CreatedAt(),
		UpdatedAt: party.UpdatedAt(),
	}
}
