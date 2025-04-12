package response

import (
	"github.com/devlucassantos/vnc-domains/src/domains/proposition"
	"github.com/google/uuid"
	"strings"
	"time"
)

type PropositionArticle struct {
	Id                   uuid.UUID        `json:"id"`
	OriginalTextUrl      string           `json:"original_text_url"`
	OriginalTextMimeType string           `json:"original_text_mime_type"`
	Title                string           `json:"title"`
	Content              string           `json:"content"`
	SubmittedAt          time.Time        `json:"submitted_at"`
	ImageUrl             string           `json:"image_url,omitempty"`
	ImageDescription     string           `json:"image_description,omitempty"`
	Deputies             []Deputy         `json:"deputies,omitempty"`
	ExternalAuthors      []ExternalAuthor `json:"external_authors,omitempty"`
	Type                 *ArticleType     `json:"type"`
	AverageRating        float64          `json:"average_rating,omitempty"`
	NumberOfRatings      int              `json:"number_of_ratings,omitempty"`
	UserRating           int              `json:"user_rating,omitempty"`
	ViewLater            bool             `json:"view_later,omitempty"`
	Votes                []Article        `json:"votes,omitempty"`
	Events               []Article        `json:"events,omitempty"`
	Newsletter           *Article         `json:"newsletter,omitempty"`
	CreatedAt            time.Time        `json:"created_at"`
	UpdatedAt            time.Time        `json:"updated_at"`
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
	articleType := NewArticleType(propositionArticle.Type())

	propositionSpecificType := proposition.Type()
	if !propositionSpecificType.IsZero() {
		articleType.SpecificType = NewPropositionSpecificType(propositionSpecificType)
	}

	var newsletter *Article
	var votes, events []Article
	for _, article := range proposition.RelatedArticles() {
		articleResponse := *NewArticle(article)
		if strings.Contains(articleResponse.Type.Codes, "voting") {
			votes = append(votes, articleResponse)
		} else if strings.Contains(articleResponse.Type.Codes, "event") {
			events = append(events, articleResponse)
		} else {
			newsletter = &articleResponse
		}
	}

	return &PropositionArticle{
		Id:                   propositionArticle.Id(),
		OriginalTextUrl:      proposition.OriginalTextUrl(),
		OriginalTextMimeType: proposition.OriginalTextMimeType(),
		Title:                proposition.Title(),
		Content:              proposition.Content(),
		SubmittedAt:          proposition.SubmittedAt(),
		ImageUrl:             proposition.ImageUrl(),
		ImageDescription:     proposition.ImageDescription(),
		Deputies:             deputies,
		ExternalAuthors:      externalAuthors,
		Type:                 articleType,
		AverageRating:        propositionArticle.AverageRating(),
		NumberOfRatings:      propositionArticle.NumberOfRatings(),
		UserRating:           propositionArticle.UserRating(),
		ViewLater:            propositionArticle.ViewLater(),
		Votes:                votes,
		Events:               events,
		Newsletter:           newsletter,
		CreatedAt:            propositionArticle.CreatedAt(),
		UpdatedAt:            propositionArticle.UpdatedAt(),
	}
}
