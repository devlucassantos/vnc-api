package repositories

import (
	"github.com/google/uuid"
	"vnc-read-api/core/domains/proposition"
)

type Proposition interface {
	GetPropositionById(id uuid.UUID) (*proposition.Proposition, error)
}
