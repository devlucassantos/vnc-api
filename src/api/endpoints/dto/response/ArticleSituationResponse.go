package response

import (
	"github.com/devlucassantos/vnc-domains/src/domains/articlesituation"
	"github.com/google/uuid"
	"time"
)

type ArticleSituation struct {
	Id          *uuid.UUID `json:"id,omitempty"`
	Description string     `json:"description,omitempty"`
	Color       string     `json:"color,omitempty"`
	StartsAt    *time.Time `json:"starts_at,omitempty"`
	EndsAt      *time.Time `json:"ends_at,omitempty"`
	IsApproved  *bool      `json:"is_approved,omitempty"`
}

func NewArticleSituation(articleSituation articlesituation.ArticleSituation) *ArticleSituation {
	var idPointer *uuid.UUID
	if articleSituation.Id() != uuid.Nil {
		id := articleSituation.Id()
		idPointer = &id
	}

	var startsAtPointer *time.Time
	if !articleSituation.StartsAt().IsZero() {
		startsAt := articleSituation.StartsAt()
		startsAtPointer = &startsAt
	}

	var endsAtPointer *time.Time
	if !articleSituation.EndsAt().IsZero() {
		endsAt := articleSituation.EndsAt()
		endsAtPointer = &endsAt
	}

	var isApprovedPointer *bool
	if articleSituation.IsApproved() {
		isApproved := articleSituation.IsApproved()
		isApprovedPointer = &isApproved
	}

	return &ArticleSituation{
		Id:          idPointer,
		Description: articleSituation.Description(),
		Color:       articleSituation.Color(),
		StartsAt:    startsAtPointer,
		EndsAt:      endsAtPointer,
		IsApproved:  isApprovedPointer,
	}
}
