package swagger

import "github.com/google/uuid"

type AgendaItemRegime struct {
	Id          uuid.UUID `json:"id"          example:"667084f2-eddb-495b-b50f-8f60a2949a92"`
	Description string    `json:"description" example:"Matéria Sobre a Mesa"`
}
