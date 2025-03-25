package swagger

import (
	"github.com/google/uuid"
	"time"
)

type VotingArticle struct {
	Id                   uuid.UUID       `json:"id"                    example:"d369b9bc-c226-4bbf-8fbb-fceed205845a"`
	Title                string          `json:"title"                 example:"Votação 3457539-42"`
	Description          string          `json:"description"           example:"A Câmara dos Deputados aprovou o Projeto de Lei nº 1..."`
	Result               string          `json:"result"                example:"Aprovado o Substitutivo ao Projeto de Lei nº 1..."`
	ResultAnnouncedAt    time.Time       `json:"result_announced_at"   example:"2023-05-18T20:17:32Z"`
	IsApproved           bool            `json:"is_approved"           example:"true"`
	LegislativeBody      LegislativeBody `json:"legislative_body"`
	MainProposition      Article         `json:"main_proposition"`
	RelatedPropositions  []Article       `json:"related_propositions"`
	AffectedPropositions []Article       `json:"affected_propositions"`
	Type                 ArticleType     `json:"type"`
	AverageRating        float64         `json:"average_rating"        example:"3.5"`
	NumberOfRatings      int             `json:"number_of_ratings"     example:"254"`
	UserRating           int             `json:"user_rating"           example:"4"`
	ViewLater            bool            `json:"view_later"            example:"true"`
	Events               []Article       `json:"events"`
	Newsletter           Article         `json:"newsletter"`
	CreatedAt            time.Time       `json:"created_at"            example:"2023-05-18T23:19:00.465814Z"`
	UpdatedAt            time.Time       `json:"updated_at"            example:"2023-05-18T23:19:00.465814Z"`
}
