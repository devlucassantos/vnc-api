package services

import (
	"github.com/devlucassantos/vnc-domains/src/domains/articletype"
	"github.com/devlucassantos/vnc-domains/src/domains/deputy"
	"github.com/devlucassantos/vnc-domains/src/domains/eventsituation"
	"github.com/devlucassantos/vnc-domains/src/domains/eventtype"
	"github.com/devlucassantos/vnc-domains/src/domains/externalauthor"
	"github.com/devlucassantos/vnc-domains/src/domains/legislativebody"
	"github.com/devlucassantos/vnc-domains/src/domains/party"
	"github.com/devlucassantos/vnc-domains/src/domains/propositiontype"
)

type Resources interface {
	GetResources() ([]articletype.ArticleType, []propositiontype.PropositionType, []party.Party, []deputy.Deputy,
		[]externalauthor.ExternalAuthor, []legislativebody.LegislativeBody, []eventtype.EventType,
		[]eventsituation.EventSituation, error)
	GetArticleTypes() ([]articletype.ArticleType, error)
	GetPropositionTypes() ([]propositiontype.PropositionType, error)
	GetParties() ([]party.Party, error)
	GetDeputies() ([]deputy.Deputy, error)
	GetExternalAuthors() ([]externalauthor.ExternalAuthor, error)
	GetLegislativeBodies() ([]legislativebody.LegislativeBody, error)
	GetEventTypes() ([]eventtype.EventType, error)
	GetEventSituations() ([]eventsituation.EventSituation, error)
}
