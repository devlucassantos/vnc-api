package response

import (
	"github.com/google/uuid"
	"time"
	"vnc-read-api/core/domains/party"
)

type Party struct {
	Id        uuid.UUID `json:"id,omitempty"`
	Code      int       `json:"code,omitempty"`
	Name      string    `json:"name,omitempty"`
	Acronym   string    `json:"acronym,omitempty"`
	ImageUrl  string    `json:"image_url,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
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
