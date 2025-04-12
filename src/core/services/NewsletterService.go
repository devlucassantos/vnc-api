package services

import (
	"github.com/devlucassantos/vnc-domains/src/domains/newsletter"
	"github.com/google/uuid"
	"vnc-api/core/interfaces/postgres"
)

type Newsletter struct {
	repository postgres.Newsletter
}

func NewNewsletterService(repository postgres.Newsletter) *Newsletter {
	return &Newsletter{
		repository: repository,
	}
}

func (instance Newsletter) GetNewsletterByArticleId(articleId uuid.UUID, userId uuid.UUID) (*newsletter.Newsletter, error) {
	return instance.repository.GetNewsletterByArticleId(articleId, userId)
}
