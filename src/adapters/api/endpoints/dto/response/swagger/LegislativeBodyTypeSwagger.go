package swagger

import "github.com/google/uuid"

type LegislativeBodyType struct {
	Id          uuid.UUID `json:"id"          example:"3440489a-d787-447f-80dd-51c0c577f07f"`
	Description string    `json:"description" example:"Comissão Mista Permanente"`
}
