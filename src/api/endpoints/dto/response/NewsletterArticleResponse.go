package response

import (
	"github.com/devlucassantos/vnc-domains/src/domains/newsletter"
	"github.com/google/uuid"
	"time"
)

type NewsletterArticle struct {
	Id                  uuid.UUID `json:"id"`
	Title               string    `json:"title"`
	Content             string    `json:"content"`
	ReferenceDate       time.Time `json:"reference_date"`
	AverageRating       float64   `json:"average_rating,omitempty"`
	NumberOfRatings     int       `json:"number_of_ratings,omitempty"`
	UserRating          int       `json:"user_rating,omitempty"`
	ViewLater           bool      `json:"view_later"`
	PropositionArticles []Article `json:"proposition_articles"`
	ReferenceDateTime   time.Time `json:"reference_date_time"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

func NewNewsletterArticle(newsletter newsletter.Newsletter) *NewsletterArticle {
	newsletterArticle := newsletter.Article()

	return &NewsletterArticle{
		Id:                newsletterArticle.Id(),
		Title:             newsletter.Title(),
		Content:           newsletter.Description(),
		ReferenceDate:     newsletter.ReferenceDate(),
		AverageRating:     newsletterArticle.AverageRating(),
		NumberOfRatings:   newsletterArticle.NumberOfRatings(),
		UserRating:        newsletterArticle.UserRating(),
		ViewLater:         newsletterArticle.ViewLater(),
		ReferenceDateTime: newsletterArticle.ReferenceDateTime(),
		CreatedAt:         newsletterArticle.CreatedAt(),
		UpdatedAt:         newsletterArticle.UpdatedAt(),
	}
}
