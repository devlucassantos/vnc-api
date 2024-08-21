package response

import (
	"github.com/devlucassantos/vnc-domains/src/domains/proptype"
	"github.com/google/uuid"
	"time"
)

type PropositionType struct {
	Id                   uuid.UUID `json:"id"`
	Description          string    `json:"description"`
	Color                string    `json:"color"`
	SortOrder            int       `json:"sort_order"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
	PropositionsArticles []Article `json:"proposition_articles,omitempty"`
}

func NewPropositionType(propositionType proptype.PropositionType) *PropositionType {
	return &PropositionType{
		Id:          propositionType.Id(),
		Description: propositionType.Description(),
		Color:       propositionType.Color(),
		SortOrder:   propositionType.SortOrder(),
		CreatedAt:   propositionType.CreatedAt(),
		UpdatedAt:   propositionType.UpdatedAt(),
	}
}
