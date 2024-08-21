package response

import (
	"github.com/devlucassantos/vnc-domains/src/domains/deputy"
	"github.com/devlucassantos/vnc-domains/src/domains/external"
	"github.com/devlucassantos/vnc-domains/src/domains/party"
	"github.com/devlucassantos/vnc-domains/src/domains/proptype"
)

type Resources struct {
	PropositionTypes []PropositionType `json:"proposition_types"`
	Parties          []Party           `json:"parties"`
	Deputies         []Deputy          `json:"deputies"`
	ExternalAuthors  []ExternalAuthor  `json:"external_authors"`
}

func NewResources(propositionTypes []proptype.PropositionType, parties []party.Party, deputies []deputy.Deputy,
	externalAuthors []external.ExternalAuthor) *Resources {
	propositionTypeList := []PropositionType{}
	for _, propositionTypeData := range propositionTypes {
		propositionTypeList = append(propositionTypeList, *NewPropositionType(propositionTypeData))
	}

	partyList := []Party{}
	for _, partyData := range parties {
		partyList = append(partyList, *NewParty(partyData))
	}

	deputyList := []Deputy{}
	for _, deputyData := range deputies {
		deputyList = append(deputyList, *NewDeputy(deputyData))
	}

	externalAuthorList := []ExternalAuthor{}
	for _, externalAuthorData := range externalAuthors {
		externalAuthorList = append(externalAuthorList, *NewExternalAuthor(externalAuthorData))
	}

	return &Resources{
		PropositionTypes: propositionTypeList,
		Parties:          partyList,
		Deputies:         deputyList,
		ExternalAuthors:  externalAuthorList,
	}
}
