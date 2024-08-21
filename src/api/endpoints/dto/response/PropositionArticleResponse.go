package response

import (
	"github.com/devlucassantos/vnc-domains/src/domains/proposition"
	"github.com/google/uuid"
	"time"
)

type PropositionArticle struct {
	Id                uuid.UUID        `json:"id"`
	OriginalTextUrl   string           `json:"original_text_url"`
	Title             string           `json:"title"`
	Content           string           `json:"content"`
	SubmittedAt       time.Time        `json:"submitted_at"`
	ImageUrl          string           `json:"image_url"`
	Type              *PropositionType `json:"type"`
	Deputies          []Deputy         `json:"deputies,omitempty"`
	ExternalAuthors   []ExternalAuthor `json:"externalAuthors,omitempty"`
	AverageRating     float64          `json:"average_rating,omitempty"`
	NumberOfRatings   int              `json:"number_of_ratings,omitempty"`
	UserRating        int              `json:"user_rating,omitempty"`
	ViewLater         bool             `json:"view_later"`
	NewsletterArticle *Article         `json:"newsletter_article,omitempty"`
	ReferenceDateTime time.Time        `json:"reference_date_time"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
}

func NewPropositionArticle(proposition proposition.Proposition) *PropositionArticle {
	var deputies []Deputy
	for _, deputy := range proposition.Deputies() {
		deputies = append(deputies, *NewDeputy(deputy))
	}

	var externalAuthors []ExternalAuthor
	for _, externalAuthor := range proposition.ExternalAuthors() {
		externalAuthors = append(externalAuthors, *NewExternalAuthor(externalAuthor))
	}

	propositionArticle := proposition.Article()

	return &PropositionArticle{
		Id:                propositionArticle.Id(),
		OriginalTextUrl:   proposition.OriginalTextUrl(),
		Title:             proposition.Title(),
		Content:           proposition.Content(),
		SubmittedAt:       proposition.SubmittedAt(),
		ImageUrl:          proposition.ImageUrl(),
		Type:              NewPropositionType(proposition.Type()),
		Deputies:          deputies,
		ExternalAuthors:   externalAuthors,
		AverageRating:     propositionArticle.AverageRating(),
		NumberOfRatings:   propositionArticle.NumberOfRatings(),
		UserRating:        propositionArticle.UserRating(),
		ViewLater:         propositionArticle.ViewLater(),
		ReferenceDateTime: propositionArticle.ReferenceDateTime(),
		CreatedAt:         propositionArticle.CreatedAt(),
		UpdatedAt:         propositionArticle.UpdatedAt(),
	}
}
