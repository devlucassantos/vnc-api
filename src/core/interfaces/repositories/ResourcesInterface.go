package repositories

import (
	"github.com/devlucassantos/vnc-domains/src/domains/articletype"
	"github.com/devlucassantos/vnc-domains/src/domains/deputy"
	"github.com/devlucassantos/vnc-domains/src/domains/external"
	"github.com/devlucassantos/vnc-domains/src/domains/party"
	"github.com/google/uuid"
)

type Resources interface {
	GetArticleTypes(articleTypeIds []uuid.UUID) ([]articletype.ArticleType, error)
	GetParties() ([]party.Party, error)
	GetDeputies() ([]deputy.Deputy, error)
	GetExternalAuthors() ([]external.ExternalAuthor, error)
}
