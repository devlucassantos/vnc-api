package swagger

import "github.com/google/uuid"

type EventSituation struct {
	Id          uuid.UUID `json:"id"          example:"0f7f6406-bbce-4a26-aa02-c1694bda5e72"`
	Description string    `json:"description" example:"Encerrado"`
	Color       string    `json:"color"       example:"#0047AB"`
}
