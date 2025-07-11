package dto

import (
	"github.com/google/uuid"
)

type EventSituation struct {
	Id          uuid.UUID `db:"event_situation_id"`
	Description string    `db:"event_situation_description"`
	Color       string    `db:"event_situation_color"`
}
