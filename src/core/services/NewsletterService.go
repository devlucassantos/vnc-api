package services

import (
	"github.com/devlucassantos/vnc-domains/src/domains/newsletter"
	"github.com/google/uuid"
	"vnc-api/core/interfaces/repositories"
)

type Newsletter struct {
	repository repositories.Newsletter
}

func NewNewsletterService(repository repositories.Newsletter) *Newsletter {
	return &Newsletter{
		repository: repository,
	}
}

func (instance Newsletter) GetNewsletterByArticleId(articleId uuid.UUID, userId uuid.UUID) (*newsletter.Newsletter, error) {
	return instance.repository.GetNewsletterByArticleId(articleId, userId)
}
