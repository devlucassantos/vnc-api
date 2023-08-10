package repositories

import (
	"github.com/google/uuid"
	"vnc-read-api/api/endpoints/dto/filter"
	"vnc-read-api/core/domains/proposition"
)

type Proposition interface {
	GetPropositions(filter filter.PropositionFilter) ([]proposition.Proposition, int, error)
	GetPropositionById(id uuid.UUID) (*proposition.Proposition, error)
}
