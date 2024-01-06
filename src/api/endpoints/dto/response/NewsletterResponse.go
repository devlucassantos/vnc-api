package response

import (
	"github.com/google/uuid"
	"time"
	"vnc-read-api/core/domains/newsletter"
)

type Newsletter struct {
	Id            uuid.UUID     `json:"id,omitempty"`
	Title         string        `json:"title,omitempty"`
	Content       string        `json:"content,omitempty"`
	ReferenceDate time.Time     `json:"reference_date,omitempty"`
	Propositions  []Proposition `json:"propositions,omitempty"`
	CreatedAt     time.Time     `json:"created_at,omitempty"`
	UpdatedAt     time.Time     `json:"updated_at,omitempty"`
}

func NewNewsletter(newsletter newsletter.Newsletter) *Newsletter {
	var propositions []Proposition
	for _, propositionData := range newsletter.Propositions() {
		propositions = append(propositions, *NewProposition(propositionData))
	}

	return &Newsletter{
		Id:            newsletter.Id(),
		Title:         newsletter.Title(),
		Content:       newsletter.Content(),
		ReferenceDate: newsletter.ReferenceDate(),
		Propositions:  propositions,
		CreatedAt:     newsletter.CreatedAt(),
		UpdatedAt:     newsletter.UpdatedAt(),
	}
}
