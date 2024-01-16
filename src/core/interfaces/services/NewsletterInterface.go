package services

import (
	"github.com/devlucassantos/vnc-domains/src/domains/newsletter"
	"github.com/google/uuid"
)

type Newsletter interface {
	GetNewsletterById(id uuid.UUID) (*newsletter.Newsletter, error)
	GetNewsletterByPropositionId(propositionId uuid.UUID) (*newsletter.Newsletter, error)
}
