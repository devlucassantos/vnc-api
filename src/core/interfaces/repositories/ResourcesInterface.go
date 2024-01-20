package repositories

import (
	"github.com/devlucassantos/vnc-domains/src/domains/deputy"
	"github.com/devlucassantos/vnc-domains/src/domains/organization"
	"github.com/devlucassantos/vnc-domains/src/domains/party"
)

type Resources interface {
	GetParties() ([]party.Party, error)
	GetDeputies() ([]deputy.Deputy, error)
	GetOrganizations() ([]organization.Organization, error)
}
