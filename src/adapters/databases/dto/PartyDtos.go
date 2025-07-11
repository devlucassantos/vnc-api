package dto

import (
	"github.com/google/uuid"
)

type Party struct {
	Id       uuid.UUID `db:"party_id"`
	Name     string    `db:"party_name"`
	Acronym  string    `db:"party_acronym"`
	ImageUrl string    `db:"party_image_url"`
}

type PreviousParty struct {
	Id       uuid.UUID `db:"previous_party_id"`
	Name     string    `db:"previous_party_name"`
	Acronym  string    `db:"previous_party_acronym"`
	ImageUrl string    `db:"previous_party_image_url"`
}
