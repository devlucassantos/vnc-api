package swagger

import (
	"github.com/google/uuid"
	"time"
)

type Article struct {
	Id                    uuid.UUID                   `json:"id"                     example:"b27947d6-3224-4479-8da4-7917ae16b34d"`
	Title                 string                      `json:"title"                  example:"Novo projeto de lei impulsiona crescimento do setor portuário até 2028"`
	Content               string                      `json:"content"                example:"Foi sancionado o projeto de lei que altera a Lei n° 11.033 para prorrogar o Regime Tributário..."`
	MultimediaUrl         string                      `json:"multimedia_url"         example:"https://image.url.com/article/b27947d6-3224-4479-8da4-7917ae16b34d/image.png"`
	MultimediaDescription string                      `json:"multimedia_description" example:"A imagem retrata uma pequena cidade costeira afetada por uma forte chuva..."`
	Situation             ArticleSituation            `json:"situation"`
	AverageRating         float64                     `json:"average_rating"         example:"3.5"`
	NumberOfRatings       int                         `json:"number_of_ratings"      example:"452"`
	UserRating            int                         `json:"user_rating"            example:"3"`
	ViewLater             bool                        `json:"view_later"             example:"true"`
	Type                  ArticleTypeWithSpecificType `json:"type"`
	CreatedAt             time.Time                   `json:"created_at"             example:"2024-01-05T20:25:19.98031Z"`
	UpdatedAt             time.Time                   `json:"updated_at"             example:"2024-01-05T20:25:19.98031Z"`
}
