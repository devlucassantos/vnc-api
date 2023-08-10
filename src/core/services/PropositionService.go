package services

import (
	"github.com/google/uuid"
	"vnc-read-api/api/endpoints/dto/filter"
	"vnc-read-api/core/domains/proposition"
	"vnc-read-api/core/interfaces/repositories"
)

type Proposition struct {
	repository repositories.Proposition
}

func NewPropositionService(repository repositories.Proposition) *Proposition {
	return &Proposition{
		repository: repository,
	}
}

func (instance Proposition) GetPropositions(filter filter.PropositionFilter) ([]proposition.Proposition, int, error) {
	return instance.repository.GetPropositions(filter)
}

func (instance Proposition) GetPropositionById(id uuid.UUID) (*proposition.Proposition, error) {
	return instance.repository.GetPropositionById(id)
}
