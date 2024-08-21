package response

import (
	"github.com/devlucassantos/vnc-domains/src/domains/article"
	"github.com/google/uuid"
	"time"
)

type Article struct {
	Id                uuid.UUID `json:"id"`
	Title             string    `json:"title"`
	Content           string    `json:"content"`
	ImageUrl          string    `json:"image_url,omitempty"`
	AverageRating     float64   `json:"average_rating,omitempty"`
	NumberOfRatings   int       `json:"number_of_ratings,omitempty"`
	UserRating        int       `json:"user_rating,omitempty"`
	ViewLater         bool      `json:"view_later"`
	Type              string    `json:"type"`
	ReferenceDateTime time.Time `json:"reference_date_time"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func NewArticle(article article.Article) *Article {
	return &Article{
		Id:                article.Id(),
		Title:             article.Title(),
		Content:           article.Content(),
		ImageUrl:          article.ImageUrl(),
		AverageRating:     article.AverageRating(),
		NumberOfRatings:   article.NumberOfRatings(),
		UserRating:        article.UserRating(),
		ViewLater:         article.ViewLater(),
		Type:              article.Type(),
		ReferenceDateTime: article.ReferenceDateTime(),
		CreatedAt:         article.CreatedAt(),
		UpdatedAt:         article.UpdatedAt(),
	}
}
