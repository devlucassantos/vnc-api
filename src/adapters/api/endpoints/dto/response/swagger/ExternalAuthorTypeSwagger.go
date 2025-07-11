package swagger

import "github.com/google/uuid"

type ExternalAuthorType struct {
	Id          uuid.UUID `json:"id"          example:"344b3941-69fb-4715-a7e5-6afc21cd3e48"`
	Description string    `json:"description" example:"Órgão Do Poder Legislativo"`
}
