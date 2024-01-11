package response

import (
	"vnc-read-api/core/domains/deputy"
	"vnc-read-api/core/domains/organization"
	"vnc-read-api/core/domains/party"
)

type Resources struct {
	Parties       []Party        `json:"parties,omitempty"`
	Deputies      []Deputy       `json:"deputies,omitempty"`
	Organizations []Organization `json:"organizations,omitempty"`
}

func NewResources(parties []party.Party, deputies []deputy.Deputy, organizations []organization.Organization) *Resources {
	var partyList []Party
	for _, partyData := range parties {
		partyList = append(partyList, *NewParty(partyData))
	}

	var deputyList []Deputy
	for _, deputyData := range deputies {
		deputyList = append(deputyList, *NewDeputy(deputyData))
	}

	var organizationList []Organization
	for _, organizationData := range organizations {
		organizationList = append(organizationList, *NewOrganization(organizationData))
	}

	return &Resources{
		Parties:       partyList,
		Deputies:      deputyList,
		Organizations: organizationList,
	}
}
