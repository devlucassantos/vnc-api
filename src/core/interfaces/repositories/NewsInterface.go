package repositories

import (
	"vnc-read-api/core/domains/news"
	"vnc-read-api/core/filters"
)

type News interface {
	GetNews(filter filters.NewsFilter) ([]news.News, int, error)
	GetTrending(filter filters.NewsFilter) ([]news.News, int, error)
}
