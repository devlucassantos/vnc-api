package services

import (
	"github.com/devlucassantos/vnc-domains/src/domains/newsletter"
	"github.com/google/uuid"
	"vnc-read-api/core/interfaces/repositories"
)

type Newsletter struct {
	repository repositories.Newsletter
}

func NewNewsletterService(repository repositories.Newsletter) *Newsletter {
	return &Newsletter{
		repository: repository,
	}
}

func (instance Newsletter) GetNewsletterById(id uuid.UUID) (*newsletter.Newsletter, error) {
	return instance.repository.GetNewsletterById(id)
}

func (instance Newsletter) GetNewsletterByPropositionId(propositionId uuid.UUID) (*newsletter.Newsletter, error) {
	return instance.repository.GetNewsletterByPropositionId(propositionId)
}
