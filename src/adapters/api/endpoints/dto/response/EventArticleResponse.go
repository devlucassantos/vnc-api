package response

import (
	"github.com/devlucassantos/vnc-domains/src/domains/event"
	"github.com/google/uuid"
	"time"
)

type EventArticle struct {
	Id                 uuid.UUID         `json:"id"`
	Title              string            `json:"title"`
	DescriptionContent string            `json:"description_content"`
	StartsAt           time.Time         `json:"starts_at"`
	EndsAt             *time.Time        `json:"ends_at,omitempty"`
	Location           string            `json:"location"`
	IsInternal         bool              `json:"is_internal"`
	VideoUrl           string            `json:"video_url,omitempty"`
	LegislativeBodies  []LegislativeBody `json:"legislative_bodies,omitempty"`
	Requirements       []Article         `json:"requirements,omitempty"`
	AgendaItems        []EventAgendaItem `json:"agenda_items,omitempty"`
	Situation          *EventSituation   `json:"situation"`
	Type               *ArticleType      `json:"type"`
	AverageRating      float64           `json:"average_rating,omitempty"`
	NumberOfRatings    int               `json:"number_of_ratings,omitempty"`
	UserRating         int               `json:"user_rating,omitempty"`
	ViewLater          bool              `json:"view_later,omitempty"`
	Newsletter         *Article          `json:"newsletter,omitempty"`
	CreatedAt          time.Time         `json:"created_at"`
	UpdatedAt          time.Time         `json:"updated_at"`
}

func NewEventArticle(event event.Event) *EventArticle {
	var endsAt *time.Time
	if !event.EndsAt().IsZero() {
		eventEndsAt := event.EndsAt()
		endsAt = &eventEndsAt
	}

	var legislativeBodies []LegislativeBody
	for _, legislativeBody := range event.LegislativeBodies() {
		legislativeBodies = append(legislativeBodies, *NewLegislativeBody(legislativeBody))
	}

	var requirements []Article
	for _, requirement := range event.Requirements() {
		requirements = append(requirements, *NewArticle(requirement.Article()))
	}

	var agendaItems []EventAgendaItem
	for _, agendaItem := range event.AgendaItems() {
		agendaItems = append(agendaItems, *NewEventAgendaItem(agendaItem))
	}

	eventArticle := event.Article()
	articleType := NewArticleType(eventArticle.Type())

	eventSpecificType := event.Type()
	if !eventSpecificType.IsZero() {
		articleType.SpecificType = NewEventSpecificType(eventSpecificType)
	}

	var newsletter *Article
	for _, article := range event.RelatedArticles() {
		newsletter = NewArticle(article)
	}

	return &EventArticle{
		Id:                 eventArticle.Id(),
		Title:              event.Title(),
		DescriptionContent: event.Description(),
		StartsAt:           event.StartsAt(),
		EndsAt:             endsAt,
		Location:           event.Location(),
		IsInternal:         event.IsInternal(),
		VideoUrl:           event.VideoUrl(),
		LegislativeBodies:  legislativeBodies,
		Requirements:       requirements,
		AgendaItems:        agendaItems,
		Situation:          NewEventSituation(event.Situation()),
		Type:               articleType,
		AverageRating:      eventArticle.AverageRating(),
		NumberOfRatings:    eventArticle.NumberOfRatings(),
		UserRating:         eventArticle.UserRating(),
		ViewLater:          eventArticle.ViewLater(),
		Newsletter:         newsletter,
		CreatedAt:          eventArticle.CreatedAt(),
		UpdatedAt:          eventArticle.UpdatedAt(),
	}
}
