package response

import (
	"github.com/google/uuid"
	"time"
	"vnc-read-api/core/domains/keyword"
)

type Keyword struct {
	Id        uuid.UUID `json:"id"`
	Keyword   string    `json:"keyword"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewKeyword(keyword keyword.Keyword) *Keyword {
	return &Keyword{
		Id:        keyword.Id(),
		Keyword:   keyword.Keyword(),
		CreatedAt: keyword.CreatedAt(),
		UpdatedAt: keyword.UpdatedAt(),
	}
}
