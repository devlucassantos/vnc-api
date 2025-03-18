package swagger

import "github.com/google/uuid"

type EventType struct {
	Id          uuid.UUID `json:"id"          example:"2371ffd1-67d0-40b2-b6ee-dde3948cbb80"`
	Description string    `json:"description" example:"Reunião Deliberativa"`
	Color       string    `json:"color"       example:"#006767"`
}
