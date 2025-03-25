package response

import (
	"github.com/devlucassantos/vnc-domains/src/domains/voting"
	"github.com/google/uuid"
	"strings"
	"time"
)

type VotingArticle struct {
	Id                   uuid.UUID        `json:"id"`
	Title                string           `json:"title"`
	Description          string           `json:"description"`
	Result               string           `json:"result"`
	ResultAnnouncedAt    time.Time        `json:"result_announced_at"`
	IsApproved           bool             `json:"is_approved"`
	LegislativeBody      *LegislativeBody `json:"legislative_body"`
	MainProposition      *Article         `json:"main_proposition,omitempty"`
	RelatedPropositions  []Article        `json:"related_propositions,omitempty"`
	AffectedPropositions []Article        `json:"affected_propositions,omitempty"`
	Type                 *ArticleType     `json:"type"`
	AverageRating        float64          `json:"average_rating,omitempty"`
	NumberOfRatings      int              `json:"number_of_ratings,omitempty"`
	UserRating           int              `json:"user_rating,omitempty"`
	ViewLater            bool             `json:"view_later,omitempty"`
	Events               []Article        `json:"events,omitempty"`
	Newsletter           *Article         `json:"newsletter,omitempty"`
	CreatedAt            time.Time        `json:"created_at"`
	UpdatedAt            time.Time        `json:"updated_at"`
}

func NewVotingArticle(voting voting.Voting) *VotingArticle {
	mainProposition := voting.MainProposition()
	var mainPropositionArticle *Article
	if mainProposition.Id() != uuid.Nil {
		mainPropositionArticle = NewArticle(mainProposition.Article())
	}

	var relatedPropositions []Article
	for _, proposition := range voting.RelatedPropositions() {
		relatedPropositions = append(relatedPropositions, *NewArticle(proposition.Article()))
	}

	var affectedPropositions []Article
	for _, proposition := range voting.AffectedPropositions() {
		affectedPropositions = append(affectedPropositions, *NewArticle(proposition.Article()))
	}

	votingArticle := voting.Article()

	var newsletter *Article
	var events []Article
	for _, article := range voting.RelatedArticles() {
		articleResponse := *NewArticle(article)
		if strings.Contains(articleResponse.Type.Codes, "event") {
			events = append(events, articleResponse)
		} else {
			newsletter = &articleResponse
		}
	}

	return &VotingArticle{
		Id:                   voting.Id(),
		Title:                voting.Title(),
		Description:          voting.Description(),
		Result:               voting.Result(),
		ResultAnnouncedAt:    voting.ResultAnnouncedAt(),
		IsApproved:           voting.IsApproved(),
		LegislativeBody:      NewLegislativeBody(voting.LegislativeBody()),
		MainProposition:      mainPropositionArticle,
		RelatedPropositions:  relatedPropositions,
		AffectedPropositions: affectedPropositions,
		Type:                 NewArticleType(votingArticle.Type()),
		AverageRating:        votingArticle.AverageRating(),
		NumberOfRatings:      votingArticle.NumberOfRatings(),
		UserRating:           votingArticle.UserRating(),
		ViewLater:            votingArticle.ViewLater(),
		Events:               events,
		Newsletter:           newsletter,
		CreatedAt:            votingArticle.CreatedAt(),
		UpdatedAt:            votingArticle.UpdatedAt(),
	}
}
