package services

import (
	"github.com/devlucassantos/vnc-domains/src/domains/deputy"
	"github.com/devlucassantos/vnc-domains/src/domains/external"
	"github.com/devlucassantos/vnc-domains/src/domains/party"
	"github.com/devlucassantos/vnc-domains/src/domains/proptype"
	"github.com/google/uuid"
)

type Resources interface {
	GetResources() ([]proptype.PropositionType, []party.Party, []deputy.Deputy, []external.ExternalAuthor, error)
	GetPropositionTypes(propositionTypeIds []uuid.UUID) ([]proptype.PropositionType, error)
}
