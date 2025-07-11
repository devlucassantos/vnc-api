package swagger

import (
	"github.com/google/uuid"
	"time"
)

type EventArticle struct {
	Id                 uuid.UUID                   `json:"id"                  example:"e36492f7-5305-4333-9aea-05dcf3689693"`
	Title              string                      `json:"title"               example:"TCU Pede Ação Urgente contra Dívida Pública..."`
	DescriptionContent string                      `json:"description_content" example:"Discussão e votação de propostas legislativas..."`
	StartsAt           time.Time                   `json:"starts_at"           example:"2023-08-24T13:00:00Z"`
	EndsAt             time.Time                   `json:"ends_at"             example:"2023-08-24T18:00:00Z"`
	Location           string                      `json:"location"            example:"Plenário da Câmara dos Deputados"`
	IsInternal         bool                        `json:"is_internal"         example:"true"`
	VideoUrl           string                      `json:"video_url"           example:"https://video.url.com/article/e36492f7-5305-4333-9aea-05dcf3689693/video.mp4"`
	LegislativeBodies  []LegislativeBody           `json:"legislative_bodies"`
	Requirements       []Article                   `json:"requirements"`
	AgendaItems        []EventAgendaItem           `json:"agenda_items"`
	Situation          EventSituation              `json:"situation"`
	Type               ArticleTypeWithSpecificType `json:"type"`
	AverageRating      float64                     `json:"average_rating"      example:"2.5"`
	NumberOfRatings    int                         `json:"number_of_ratings"   example:"2"`
	UserRating         int                         `json:"user_rating"         example:"4"`
	ViewLater          bool                        `json:"view_later"          example:"false"`
	Newsletter         Article                     `json:"newsletter"`
	CreatedAt          time.Time                   `json:"created_at"          example:"2023-08-22T10:39:14.465814Z"`
	UpdatedAt          time.Time                   `json:"updated_at"          example:"2023-08-22T10:39:14.465814Z"`
}
