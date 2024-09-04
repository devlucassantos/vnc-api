package services

import (
	"github.com/devlucassantos/vnc-domains/src/domains/articletype"
	"github.com/devlucassantos/vnc-domains/src/domains/deputy"
	"github.com/devlucassantos/vnc-domains/src/domains/external"
	"github.com/devlucassantos/vnc-domains/src/domains/party"
	"github.com/google/uuid"
)

type Resources interface {
	GetResources() ([]articletype.ArticleType, []party.Party, []deputy.Deputy, []external.ExternalAuthor, error)
	GetArticleTypes(articleTypeIds []uuid.UUID) ([]articletype.ArticleType, error)
}
