package services

import (
	"vnc-read-api/core/domains/news"
	"vnc-read-api/core/filters"
	"vnc-read-api/core/interfaces/repositories"
)

type News struct {
	repository repositories.News
}

func NewNewsService(repository repositories.News) *News {
	return &News{
		repository: repository,
	}
}

func (instance News) GetNews(filter filters.NewsFilter) ([]news.News, int, error) {
	return instance.repository.GetNews(filter)
}

func (instance News) GetTrending(filter filters.NewsFilter) ([]news.News, int, error) {
	return instance.repository.GetTrending(filter)
}
