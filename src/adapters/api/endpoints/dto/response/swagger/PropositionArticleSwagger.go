package swagger

import (
	"github.com/google/uuid"
	"time"
)

type PropositionArticle struct {
	Id                   uuid.UUID                   `json:"id"                      example:"9dc67bd9-674f-4e4d-9536-07485335c362"`
	OriginalTextUrl      string                      `json:"original_text_url"       example:"https://www.camara.leg.br/proposicoesWeb/prop_mostrarintegra?codteor=4865485"`
	OriginalTextMimeType string                      `json:"original_text_mime_type" example:"https://www.camara.leg.br/proposicoesWeb/prop_mostrarintegra?codteor=4865485"`
	Title                string                      `json:"title"                   example:"Requerimento de Votação Nominal-Destaque de Emenda"`
	Content              string                      `json:"content"                 example:"O presente requerimento foi elaborado pelos deputados..."`
	SubmittedAt          time.Time                   `json:"submitted_at"            example:"2023-08-09T14:25:00Z"`
	ImageUrl             string                      `json:"image_url"               example:"https://www.vnc.com.br/news/proposition/image/87624.jpg"`
	ImageDescription     string                      `json:"image_description"       example:"A imagem retrata uma pequena cidade costeira afetada por uma forte chuva..."`
	Deputies             []Deputy                    `json:"deputies"`
	ExternalAuthors      []ExternalAuthor            `json:"external_authors"`
	Type                 ArticleTypeWithSpecificType `json:"type"`
	AverageRating        float64                     `json:"average_rating"          example:"2.5"`
	NumberOfRatings      int                         `json:"number_of_ratings"       example:"249"`
	UserRating           int                         `json:"user_rating"             example:"5"`
	ViewLater            bool                        `json:"view_later"              example:"true"`
	Votes                []Article                   `json:"votes"`
	Events               []Article                   `json:"events"`
	Newsletter           Article                     `json:"newsletter"`
	CreatedAt            time.Time                   `json:"created_at"              example:"2023-08-09T14:55:00.465814Z"`
	UpdatedAt            time.Time                   `json:"updated_at"              example:"2023-08-09T14:55:00.465814Z"`
}
