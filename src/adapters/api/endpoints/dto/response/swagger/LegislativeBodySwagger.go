package swagger

import "github.com/google/uuid"

type LegislativeBody struct {
	Id      uuid.UUID           `json:"id"      example:"31991060-fa3d-4647-aea5-b17b41de6a36"`
	Name    string              `json:"name"    example:"Comissão de Educação"`
	Acronym string              `json:"acronym" example:"2009"`
	Type    LegislativeBodyType `json:"type"`
}
