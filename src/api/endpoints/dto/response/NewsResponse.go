package response

import (
	"github.com/devlucassantos/vnc-domains/src/domains/news"
	"github.com/google/uuid"
	"time"
)

type News struct {
	Id        uuid.UUID `json:"id,omitempty"`
	Title     string    `json:"title,omitempty"`
	Content   string    `json:"content,omitempty"`
	Type      string    `json:"type,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func NewNews(news news.News) *News {
	return &News{
		Id:        news.Id(),
		Title:     news.Title(),
		Content:   news.Content(),
		Type:      news.Type(),
		CreatedAt: news.CreatedAt(),
		UpdatedAt: news.UpdatedAt(),
	}
}
