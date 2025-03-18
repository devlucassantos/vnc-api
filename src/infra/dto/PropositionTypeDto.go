package dto

import (
	"github.com/google/uuid"
)

type PropositionType struct {
	Id          uuid.UUID `db:"proposition_type_id"`
	Description string    `db:"proposition_type_description"`
	Color       string    `db:"proposition_type_color"`
}
