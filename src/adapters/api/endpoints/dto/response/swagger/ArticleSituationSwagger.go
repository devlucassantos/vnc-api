package swagger

import (
	"github.com/google/uuid"
	"time"
)

type ArticleSituation struct {
	Id                uuid.UUID `json:"id"                  example:"16f6b271-fc96-4143-972c-0a2ac149dc55"`
	Description       string    `json:"description"         example:"Em Andamento"`
	Color             string    `json:"color"               example:"#D2D2D2"`
	StartsAt          time.Time `json:"starts_at"           example:"2023-02-12T10:00:00Z"`
	EndsAt            time.Time `json:"ends_at"             example:"2023-02-15T17:00:00Z"`
	Result            string    `json:"result"              example:"Aprovado o Substitutivo ao Projeto de Lei nº 1..."`
	ResultAnnouncedAt time.Time `json:"result_announced_at" example:"2023-01-18T20:17:32Z"`
	IsApproved        bool      `json:"is_approved"         example:"true"`
}
