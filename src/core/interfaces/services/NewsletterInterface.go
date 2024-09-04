package services

import (
	"github.com/devlucassantos/vnc-domains/src/domains/newsletter"
	"github.com/google/uuid"
)

type Newsletter interface {
	GetNewsletterByArticleId(articleId uuid.UUID, userId uuid.UUID) (*newsletter.Newsletter, error)
}
