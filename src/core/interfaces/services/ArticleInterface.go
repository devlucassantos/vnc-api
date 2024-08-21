package services

import (
	"github.com/devlucassantos/vnc-domains/src/domains/article"
	"github.com/google/uuid"
	"vnc-api/core/filters"
)

type Article interface {
	GetArticles(filter filters.ArticleFilter, userId uuid.UUID) ([]article.Article, int, error)
	GetTrendingArticles(filter filters.ArticleFilter, userId uuid.UUID) ([]article.Article, int, error)
	GetTrendingArticlesByPropositionType(propositionTypeId uuid.UUID, itemsPerType int, userId uuid.UUID) ([]article.Article, error)
	GetArticlesToViewLater(filter filters.ArticleFilter, userId uuid.UUID) ([]article.Article, int, error)
	GetPropositionArticlesByNewsletterId(newsletterId uuid.UUID, userId uuid.UUID) ([]article.Article, error)
	GetNewsletterArticleByPropositionId(propositionId uuid.UUID, userId uuid.UUID) (*article.Article, error)
	SaveArticleRating(userId uuid.UUID, articleId uuid.UUID, rating int) error
	SaveArticleToViewLater(userId uuid.UUID, articleId uuid.UUID, viewLater bool) error
}
