package dto

import (
	"github.com/google/uuid"
)

type EventAgendaItem struct {
	Id                          uuid.UUID `db:"event_agenda_item_id"`
	Title                       string    `db:"event_agenda_item_title"`
	Topic                       string    `db:"event_agenda_item_topic"`
	Situation                   string    `db:"event_agenda_item_situation"`
	PropositionArticleId        uuid.UUID `db:"event_agenda_item_proposition_article_id"`
	RelatedPropositionArticleId uuid.UUID `db:"event_agenda_item_related_proposition_article_id"`
	VotingArticleId             uuid.UUID `db:"event_agenda_item_voting_article_id"`
	*AgendaItemRegime
	*Deputy
}
