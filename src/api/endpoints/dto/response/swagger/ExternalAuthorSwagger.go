package swagger

import "github.com/google/uuid"

type ExternalAuthor struct {
	Id   uuid.UUID          `json:"id"         example:"9d543e6f-20e3-4895-83e4-26b6a976580e"`
	Name string             `json:"name"       example:"Organização Você na Câmara"`
	Type ExternalAuthorType `json:"type"`
}
