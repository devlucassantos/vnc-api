package swagger

import "github.com/google/uuid"

type PropositionType struct {
	Id          uuid.UUID `json:"id"          example:"111c1a6d-d061-40b2-ad39-ec714f05c81c"`
	Description string    `json:"description" example:"Projeto de Lei"`
	Color       string    `json:"color"       example:"#C4170C"`
}
