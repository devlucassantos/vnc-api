package swagger

import "github.com/google/uuid"

type EventAgendaItem struct {
	Id                 uuid.UUID        `json:"id"        example:"a1b5ae2e-b3e7-443d-b87d-dca789224578"`
	Title              string           `json:"title"     example:"PL 132/4567"`
	Topic              string           `json:"topic"     example:"Discussão"`
	Situation          string           `json:"situation" example:"Discussão em turno único. Encerrada a discussão..."`
	Regime             AgendaItemRegime `json:"regime"`
	Rapporteur         Deputy           `json:"rapporteur"`
	Proposition        Article          `json:"proposition"`
	RelatedProposition Article          `json:"related_proposition"`
	Voting             Article          `json:"voting"`
}
