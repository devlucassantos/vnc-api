package repositories

import (
	"github.com/devlucassantos/vnc-domains/src/domains/deputy"
	"github.com/devlucassantos/vnc-domains/src/domains/external"
	"github.com/devlucassantos/vnc-domains/src/domains/party"
	"github.com/devlucassantos/vnc-domains/src/domains/proptype"
	"github.com/google/uuid"
)

type Resources interface {
	GetPropositionTypes(propositionTypeIds []uuid.UUID) ([]proptype.PropositionType, error)
	GetParties() ([]party.Party, error)
	GetDeputies() ([]deputy.Deputy, error)
	GetExternalAuthors() ([]external.ExternalAuthor, error)
}
