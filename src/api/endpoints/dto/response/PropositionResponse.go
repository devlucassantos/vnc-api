package response

import (
	"github.com/google/uuid"
	"time"
	"vnc-read-api/core/domains/proposition"
)

type Proposition struct {
	Id              uuid.UUID      `json:"id"`
	Code            int            `json:"code"`
	OriginalTextUrl string         `json:"original_text_url"`
	Title           string         `json:"title"`
	Summary         string         `json:"summary"`
	SubmittedAt     time.Time      `json:"submitted_at"`
	Deputies        []Deputy       `json:"deputies"`
	Organizations   []Organization `json:"organizations"`
	Keywords        []Keyword      `json:"keywords"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

func NewProposition(proposition proposition.Proposition) Proposition {
	var deputies []Deputy
	for _, deputyData := range proposition.Deputies() {
		deputies = append(deputies, *NewDeputy(deputyData))
	}

	var organizations []Organization
	for _, organizationData := range proposition.Organizations() {
		organizations = append(organizations, *NewOrganization(organizationData))
	}

	var keywords []Keyword
	for _, keyword := range proposition.Keywords() {
		keywords = append(keywords, *NewKeyword(keyword))
	}

	return Proposition{
		Id:              proposition.Id(),
		Code:            proposition.Code(),
		OriginalTextUrl: proposition.OriginalTextUrl(),
		Title:           proposition.Title(),
		Summary:         proposition.Summary(),
		SubmittedAt:     proposition.SubmittedAt(),
		Deputies:        deputies,
		Organizations:   organizations,
		Keywords:        keywords,
		CreatedAt:       proposition.CreatedAt(),
		UpdatedAt:       proposition.UpdatedAt(),
	}
}
