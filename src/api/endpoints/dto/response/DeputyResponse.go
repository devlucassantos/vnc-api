package response

import (
	"github.com/google/uuid"
	"time"
	"vnc-read-api/core/domains/deputy"
)

type Deputy struct {
	Id            uuid.UUID `json:"id"`
	Code          int       `json:"code"`
	Cpf           string    `json:"cpf"`
	Name          string    `json:"name"`
	ElectoralName string    `json:"electoral_name"`
	ImageUrl      string    `json:"image_url"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Party         *Party    `json:"party"`
}

func NewDeputy(deputy deputy.Deputy) *Deputy {
	return &Deputy{
		Id:            deputy.Id(),
		Code:          deputy.Code(),
		Cpf:           deputy.Cpf(),
		Name:          deputy.Name(),
		ElectoralName: deputy.ElectoralName(),
		ImageUrl:      deputy.ImageUrl(),
		CreatedAt:     deputy.CreatedAt(),
		UpdatedAt:     deputy.CreatedAt(),
		Party:         NewParty(deputy.CurrentParty()),
	}
}
