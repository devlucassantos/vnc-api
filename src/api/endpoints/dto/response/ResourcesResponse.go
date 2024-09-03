package response

import (
	"github.com/devlucassantos/vnc-domains/src/domains/articletype"
	"github.com/devlucassantos/vnc-domains/src/domains/deputy"
	"github.com/devlucassantos/vnc-domains/src/domains/external"
	"github.com/devlucassantos/vnc-domains/src/domains/party"
)

type Resources struct {
	ArticleTypes    []ArticleType    `json:"article_types"`
	Parties         []Party          `json:"parties"`
	Deputies        []Deputy         `json:"deputies"`
	ExternalAuthors []ExternalAuthor `json:"external_authors"`
}

func NewResources(articleTypes []articletype.ArticleType, parties []party.Party, deputies []deputy.Deputy,
	externalAuthors []external.ExternalAuthor) *Resources {
	articleTypeSlice := []ArticleType{}
	for _, articleTypeData := range articleTypes {
		articleTypeSlice = append(articleTypeSlice, *NewArticleType(articleTypeData))
	}

	partySlice := []Party{}
	for _, partyData := range parties {
		partySlice = append(partySlice, *NewParty(partyData))
	}

	deputySlice := []Deputy{}
	for _, deputyData := range deputies {
		deputySlice = append(deputySlice, *NewDeputy(deputyData))
	}

	externalAuthorSlice := []ExternalAuthor{}
	for _, externalAuthorData := range externalAuthors {
		externalAuthorSlice = append(externalAuthorSlice, *NewExternalAuthor(externalAuthorData))
	}

	return &Resources{
		ArticleTypes:    articleTypeSlice,
		Parties:         partySlice,
		Deputies:        deputySlice,
		ExternalAuthors: externalAuthorSlice,
	}
}
