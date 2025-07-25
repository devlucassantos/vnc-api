package dto

import (
	"github.com/google/uuid"
)

type EventType struct {
	Id          uuid.UUID `db:"event_type_id"`
	Description string    `db:"event_type_description"`
	Color       string    `db:"event_type_color"`
}
