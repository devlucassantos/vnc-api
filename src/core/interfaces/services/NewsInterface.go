package services

import (
	"vnc-read-api/api/endpoints/dto/filter"
	"vnc-read-api/core/domains/news"
)

type News interface {
	GetNews(filter filter.NewsFilter) ([]news.News, int, error)
}
