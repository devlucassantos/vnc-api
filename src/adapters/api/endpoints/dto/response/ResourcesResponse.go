package response

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

type Resources struct {
	ArticleTypes      []ArticleType     `json:"article_types"`
	PropositionTypes  []PropositionType `json:"proposition_types"`
	Parties           []Party           `json:"parties"`
	Deputies          []Deputy          `json:"deputies"`
	ExternalAuthors   []ExternalAuthor  `json:"external_authors"`
	LegislativeBodies []LegislativeBody `json:"legislative_bodies"`
	EventTypes        []EventType       `json:"event_types"`
	EventSituations   []EventSituation  `json:"event_situations"`
}

func NewResources(articleTypes []articletype.ArticleType, propositionTypes []propositiontype.PropositionType,
	parties []party.Party, deputies []deputy.Deputy, externalAuthors []externalauthor.ExternalAuthor,
	legislativeBodies []legislativebody.LegislativeBody, eventTypes []eventtype.EventType,
	eventSituations []eventsituation.EventSituation) *Resources {
	articleTypeSlice := []ArticleType{}
	for _, articleTypeData := range articleTypes {
		articleTypeSlice = append(articleTypeSlice, *NewArticleType(articleTypeData))
	}

	propositionTypeSlice := []PropositionType{}
	for _, propositionTypeData := range propositionTypes {
		propositionTypeSlice = append(propositionTypeSlice, *NewPropositionType(propositionTypeData))
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

	legislativeBodySlice := []LegislativeBody{}
	for _, legislativeBodyData := range legislativeBodies {
		legislativeBodySlice = append(legislativeBodySlice, *NewLegislativeBody(legislativeBodyData))
	}

	eventTypeSlice := []EventType{}
	for _, eventTypeData := range eventTypes {
		eventTypeSlice = append(eventTypeSlice, *NewEventType(eventTypeData))
	}

	eventSituationSlice := []EventSituation{}
	for _, eventSituationData := range eventSituations {
		eventSituationSlice = append(eventSituationSlice, *NewEventSituation(eventSituationData))
	}

	return &Resources{
		ArticleTypes:      articleTypeSlice,
		PropositionTypes:  propositionTypeSlice,
		Parties:           partySlice,
		Deputies:          deputySlice,
		ExternalAuthors:   externalAuthorSlice,
		LegislativeBodies: legislativeBodySlice,
		EventTypes:        eventTypeSlice,
		EventSituations:   eventSituationSlice,
	}
}
