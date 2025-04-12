package services

import (
	"github.com/devlucassantos/vnc-domains/src/domains/voting"
	"github.com/google/uuid"
	"vnc-api/core/interfaces/postgres"
)

type Voting struct {
	repository postgres.Voting
}

func NewVotingService(repository postgres.Voting) *Voting {
	return &Voting{
		repository: repository,
	}
}

func (instance Voting) GetVotingByArticleId(articleId uuid.UUID, userId uuid.UUID) (*voting.Voting, error) {
	return instance.repository.GetVotingByArticleId(articleId, userId)
}
