package repositories

import (
	"github.com/devlucassantos/vnc-domains/src/domains/news"
	"vnc-read-api/core/filters"
)

type News interface {
	GetNews(filter filters.NewsFilter) ([]news.News, int, error)
	GetTrending(filter filters.NewsFilter) ([]news.News, int, error)
}
