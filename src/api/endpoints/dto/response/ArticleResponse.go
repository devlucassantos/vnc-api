package response

import (
	"github.com/devlucassantos/vnc-domains/src/domains/article"
	"github.com/google/uuid"
	"strings"
	"time"
)

type Article struct {
	Id                    uuid.UUID         `json:"id"`
	Title                 string            `json:"title"`
	Content               string            `json:"content"`
	MultimediaUrl         string            `json:"multimedia_url,omitempty"`
	MultimediaDescription string            `json:"multimedia_description,omitempty"`
	Situation             *ArticleSituation `json:"situation,omitempty"`
	AverageRating         float64           `json:"average_rating,omitempty"`
	NumberOfRatings       int               `json:"number_of_ratings,omitempty"`
	UserRating            int               `json:"user_rating,omitempty"`
	ViewLater             bool              `json:"view_later,omitempty"`
	Type                  *ArticleType      `json:"type"`
	CreatedAt             time.Time         `json:"created_at"`
	UpdatedAt             time.Time         `json:"updated_at"`
}

func NewArticle(article article.Article) *Article {
	var situation *ArticleSituation
	articleSituation := article.Situation()
	if !articleSituation.IsZero() {
		situation = NewArticleSituation(articleSituation)
	}

	articleType := NewArticleType(article.Type())
	if strings.Contains(articleType.Codes, "voting") {
		isApproved := articleSituation.IsApproved()
		situation.IsApproved = &isApproved
	}

	articleSpecificType := article.SpecificType()
	if !articleSpecificType.IsZero() {
		articleType.SpecificType = NewArticleType(articleSpecificType)
	}

	return &Article{
		Id:                    article.Id(),
		Title:                 article.Title(),
		Content:               article.Content(),
		MultimediaUrl:         article.MultimediaUrl(),
		MultimediaDescription: article.MultimediaDescription(),
		Situation:             situation,
		AverageRating:         article.AverageRating(),
		NumberOfRatings:       article.NumberOfRatings(),
		UserRating:            article.UserRating(),
		ViewLater:             article.ViewLater(),
		Type:                  articleType,
		CreatedAt:             article.CreatedAt(),
		UpdatedAt:             article.UpdatedAt(),
	}
}
