package services

import (
	"github.com/devlucassantos/vnc-domains/src/domains/voting"
	"github.com/google/uuid"
)

type Voting interface {
	GetVotingByArticleId(articleId uuid.UUID, userId uuid.UUID) (*voting.Voting, error)
}
