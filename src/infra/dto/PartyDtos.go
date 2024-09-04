package dto

import (
	"github.com/google/uuid"
	"time"
)

type Party struct {
	Id        uuid.UUID `db:"party_id"`
	Name      string    `db:"party_name"`
	Acronym   string    `db:"party_acronym"`
	ImageUrl  string    `db:"party_image_url"`
	CreatedAt time.Time `db:"party_created_at"`
	UpdatedAt time.Time `db:"party_updated_at"`
}

type PartyInTheProposition struct {
	Id        uuid.UUID `db:"party_in_the_proposal_id"`
	Name      string    `db:"party_in_the_proposal_name"`
	Acronym   string    `db:"party_in_the_proposal_acronym"`
	ImageUrl  string    `db:"party_in_the_proposal_image_url"`
	CreatedAt time.Time `db:"party_in_the_proposal_created_at"`
	UpdatedAt time.Time `db:"party_in_the_proposal_updated_at"`
}
