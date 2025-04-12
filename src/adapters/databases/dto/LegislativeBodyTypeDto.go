package dto

import (
	"github.com/google/uuid"
)

type LegislativeBodyType struct {
	Id          uuid.UUID `db:"legislative_body_type_id"`
	Description string    `db:"legislative_body_type_description"`
}
