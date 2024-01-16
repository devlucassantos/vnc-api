package services

import (
	"github.com/devlucassantos/vnc-domains/src/domains/proposition"
	"github.com/google/uuid"
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

func (instance Proposition) GetPropositionById(id uuid.UUID) (*proposition.Proposition, error) {
	return instance.repository.GetPropositionById(id)
}
