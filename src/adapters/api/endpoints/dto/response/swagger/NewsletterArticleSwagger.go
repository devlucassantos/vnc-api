package swagger

import (
	"github.com/google/uuid"
	"time"
)

type NewsletterArticle struct {
	Id              uuid.UUID   `json:"id"                example:"7963a6dd-f0b8-4065-8e56-6d2a79626db7"`
	Title           string      `json:"title"             example:"Boletim do dia 26/02/2024"`
	Content         string      `json:"content"           example:"O presidente enviou ao Congresso Nacional um projeto de lei que permite a concessão de descontos fiscais..."`
	ReferenceDate   time.Time   `json:"reference_date"    example:"2023-12-23T16:34:14.441877Z"`
	Type            ArticleType `json:"type"`
	AverageRating   float64     `json:"average_rating"    example:"4.5"`
	NumberOfRatings int         `json:"number_of_ratings" example:"743"`
	UserRating      int         `json:"user_rating"       example:"5"`
	ViewLater       bool        `json:"view_later"        example:"true"`
	Propositions    []Article   `json:"propositions"`
	Votes           []Article   `json:"votes"`
	Events          []Article   `json:"events"`
	CreatedAt       time.Time   `json:"created_at"        example:"2023-12-24T19:15:22.90905Z"`
	UpdatedAt       time.Time   `json:"updated_at"        example:"2023-12-24T19:15:22.90905Z"`
}
