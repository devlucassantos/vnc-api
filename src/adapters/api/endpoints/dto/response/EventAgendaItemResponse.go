package response

import (
	"github.com/devlucassantos/vnc-domains/src/domains/eventagendaitem"
	"github.com/google/uuid"
)

type EventAgendaItem struct {
	Id                 uuid.UUID         `json:"id"`
	Title              string            `json:"title"`
	Topic              string            `json:"topic"`
	Situation          string            `json:"situation,omitempty"`
	Regime             *AgendaItemRegime `json:"regime"`
	Rapporteur         *Deputy           `json:"rapporteur,omitempty"`
	Proposition        *Article          `json:"proposition"`
	RelatedProposition *Article          `json:"related_proposition,omitempty"`
	Voting             *Article          `json:"voting,omitempty"`
}

func NewEventAgendaItem(eventAgendaItem eventagendaitem.EventAgendaItem) *EventAgendaItem {
	var rapporteur *Deputy
	eventAgendaItemRapporteur := eventAgendaItem.Rapporteur()
	if !eventAgendaItemRapporteur.IsZero() {
		rapporteur = NewDeputy(eventAgendaItemRapporteur)
	}

	proposition := eventAgendaItem.Proposition()

	var relatedPropositionArticle *Article
	relatedProposition := eventAgendaItem.RelatedProposition()
	if !relatedProposition.IsZero() {
		relatedPropositionArticle = NewArticle(relatedProposition.Article())
	}

	var votingArticle *Article
	voting := eventAgendaItem.Voting()
	if !voting.IsZero() {
		votingArticle = NewArticle(voting.Article())
	}

	return &EventAgendaItem{
		Id:                 eventAgendaItem.Id(),
		Title:              eventAgendaItem.Title(),
		Topic:              eventAgendaItem.Topic(),
		Situation:          eventAgendaItem.Situation(),
		Regime:             NewAgendaItemRegime(eventAgendaItem.Regime()),
		Rapporteur:         rapporteur,
		Proposition:        NewArticle(proposition.Article()),
		RelatedProposition: relatedPropositionArticle,
		Voting:             votingArticle,
	}
}
