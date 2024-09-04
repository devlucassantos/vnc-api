package response

import (
	"github.com/devlucassantos/vnc-domains/src/domains/articletype"
	"github.com/google/uuid"
	"time"
)

type ArticleType struct {
	Id          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	Color       string    `json:"color"`
	SortOrder   int       `json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Articles    []Article `json:"articles,omitempty"`
}

func NewArticleType(articleType articletype.ArticleType) *ArticleType {
	return &ArticleType{
		Id:          articleType.Id(),
		Description: articleType.Description(),
		Color:       articleType.Color(),
		SortOrder:   articleType.SortOrder(),
		CreatedAt:   articleType.CreatedAt(),
		UpdatedAt:   articleType.UpdatedAt(),
	}
}
