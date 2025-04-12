package services

import (
	"github.com/devlucassantos/vnc-domains/src/domains/proposition"
	"github.com/google/uuid"
	"vnc-api/core/interfaces/postgres"
)

type Proposition struct {
	repository postgres.Proposition
}

func NewPropositionService(repository postgres.Proposition) *Proposition {
	return &Proposition{
		repository: repository,
	}
}

func (instance Proposition) GetPropositionByArticleId(articleId uuid.UUID, userId uuid.UUID) (*proposition.Proposition, error) {
	return instance.repository.GetPropositionByArticleId(articleId, userId)
}
