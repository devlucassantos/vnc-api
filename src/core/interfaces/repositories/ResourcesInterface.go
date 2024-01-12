package repositories

import (
	"vnc-read-api/core/domains/deputy"
	"vnc-read-api/core/domains/organization"
	"vnc-read-api/core/domains/party"
)

type Resources interface {
	GetParties() ([]party.Party, error)
	GetDeputies() ([]deputy.Deputy, error)
	GetOrganizations() ([]organization.Organization, error)
}
