package repositories

import (
	"github.com/google/uuid"
	"vnc-read-api/core/domains/newsletter"
)

type Newsletter interface {
	GetNewsletterById(id uuid.UUID) (*newsletter.Newsletter, error)
	GetNewsletterByPropositionId(propositionId uuid.UUID) (*newsletter.Newsletter, error)
}
