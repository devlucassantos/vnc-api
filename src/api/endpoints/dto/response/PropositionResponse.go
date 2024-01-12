package response

import (
	"github.com/google/uuid"
	"time"
	"vnc-read-api/core/domains/proposition"
)

type Proposition struct {
	Id              uuid.UUID      `json:"id,omitempty"`
	Code            int            `json:"code,omitempty"`
	OriginalTextUrl string         `json:"original_text_url,omitempty"`
	Title           string         `json:"title,omitempty"`
	Content         string         `json:"content,omitempty"`
	SubmittedAt     time.Time      `json:"submitted_at,omitempty"`
	ImageUrl        string         `json:"image_url,omitempty"`
	Deputies        []Deputy       `json:"deputies,omitempty"`
	Organizations   []Organization `json:"organizations,omitempty"`
	Newsletter      *Newsletter    `json:"newsletter,omitempty"`
	CreatedAt       time.Time      `json:"created_at,omitempty"`
	UpdatedAt       time.Time      `json:"updated_at,omitempty"`
}

func NewProposition(proposition proposition.Proposition) *Proposition {
	var deputies []Deputy
	for _, deputyData := range proposition.Deputies() {
		deputies = append(deputies, *NewDeputy(deputyData))
	}

	var organizations []Organization
	for _, organizationData := range proposition.Organizations() {
		organizations = append(organizations, *NewOrganization(organizationData))
	}

	return &Proposition{
		Id:              proposition.Id(),
		Code:            proposition.Code(),
		OriginalTextUrl: proposition.OriginalTextUrl(),
		Title:           proposition.Title(),
		Content:         proposition.Content(),
		SubmittedAt:     proposition.SubmittedAt(),
		ImageUrl:        proposition.ImageUrl(),
		Deputies:        deputies,
		Organizations:   organizations,
		CreatedAt:       proposition.CreatedAt(),
		UpdatedAt:       proposition.UpdatedAt(),
	}
}
