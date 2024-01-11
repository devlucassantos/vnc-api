package services

import (
	"vnc-read-api/core/domains/deputy"
	"vnc-read-api/core/domains/organization"
	"vnc-read-api/core/domains/party"
)

type Resources interface {
	GetResources() ([]party.Party, []deputy.Deputy, []organization.Organization, error)
}
