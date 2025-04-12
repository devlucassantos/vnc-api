package response

import (
	"github.com/devlucassantos/vnc-domains/src/domains/newsletter"
	"github.com/google/uuid"
	"strings"
	"time"
)

type NewsletterArticle struct {
	Id              uuid.UUID    `json:"id"`
	Title           string       `json:"title"`
	Content         string       `json:"content"`
	Type            *ArticleType `json:"type"`
	AverageRating   float64      `json:"average_rating,omitempty"`
	NumberOfRatings int          `json:"number_of_ratings,omitempty"`
	UserRating      int          `json:"user_rating,omitempty"`
	ViewLater       bool         `json:"view_later,omitempty"`
	Propositions    []Article    `json:"propositions,omitempty"`
	Votes           []Article    `json:"votes,omitempty"`
	Events          []Article    `json:"events,omitempty"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
}

func NewNewsletterArticle(newsletter newsletter.Newsletter) *NewsletterArticle {
	newsletterArticle := newsletter.Article()

	var propositions, votes, events []Article
	for _, article := range newsletter.Articles() {
		articleResponse := *NewArticle(article)
		if strings.Contains(articleResponse.Type.Codes, "proposition") {
			propositions = append(propositions, articleResponse)
		} else if strings.Contains(articleResponse.Type.Codes, "voting") {
			votes = append(votes, articleResponse)
		} else {
			events = append(events, articleResponse)
		}
	}

	return &NewsletterArticle{
		Id:              newsletterArticle.Id(),
		Title:           newsletter.Title(),
		Content:         newsletter.Description(),
		Type:            NewArticleType(newsletterArticle.Type()),
		AverageRating:   newsletterArticle.AverageRating(),
		NumberOfRatings: newsletterArticle.NumberOfRatings(),
		UserRating:      newsletterArticle.UserRating(),
		ViewLater:       newsletterArticle.ViewLater(),
		Propositions:    propositions,
		Votes:           votes,
		Events:          events,
		CreatedAt:       newsletterArticle.CreatedAt(),
		UpdatedAt:       newsletterArticle.UpdatedAt(),
	}
}
