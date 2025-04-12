package swagger

import "github.com/google/uuid"

type ArticleType struct {
	Id          uuid.UUID `json:"id"          example:"560206f4-7360-4e21-8e45-33026f7e0953"`
	Description string    `json:"description" example:"Proposições"`
	Codes       string    `json:"codes"       example:"proposition"`
	Color       string    `json:"color"       example:"#06D13C"`
}
