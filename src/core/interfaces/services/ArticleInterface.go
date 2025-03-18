package services

import (
	"github.com/devlucassantos/vnc-domains/src/domains/article"
	"github.com/google/uuid"
	"vnc-api/core/filters"
)

type Article interface {
	GetArticles(filter filters.Article, userId uuid.UUID) ([]article.Article, int, error)
	GetTrendingArticles(filter filters.Article, userId uuid.UUID) ([]article.Article, int, error)
	GetTrendingArticlesByTypeId(articleTypeId uuid.UUID, itemsPerType int, userId uuid.UUID) ([]article.Article,
		error)
	GetTrendingArticlesBySpecificTypeId(articleSpecificTypeId uuid.UUID, itemsPerType int, userId uuid.UUID) (
		[]article.Article, error)
	GetArticlesToViewLater(filter filters.Article, userId uuid.UUID) ([]article.Article, int, error)
	SaveArticleRating(userId uuid.UUID, articleId uuid.UUID, rating *int) error
	SaveArticleToViewLater(userId uuid.UUID, articleId uuid.UUID, viewLater bool) error
}
