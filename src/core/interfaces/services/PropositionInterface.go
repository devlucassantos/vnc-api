package services

import (
	"github.com/devlucassantos/vnc-domains/src/domains/proposition"
	"github.com/google/uuid"
)

type Proposition interface {
	GetPropositionById(id uuid.UUID) (*proposition.Proposition, error)
}
