package services

import (
	"github.com/devlucassantos/vnc-domains/src/domains/deputy"
	"github.com/devlucassantos/vnc-domains/src/domains/organization"
	"github.com/devlucassantos/vnc-domains/src/domains/party"
)

type Resources interface {
	GetResources() ([]party.Party, []deputy.Deputy, []organization.Organization, error)
}
