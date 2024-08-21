package services

import (
	"github.com/devlucassantos/vnc-domains/src/domains/article"
	"github.com/google/uuid"
	"vnc-api/core/filters"
	"vnc-api/core/interfaces/repositories"
)

type Article struct {
	repository repositories.Article
}

func NewArticleService(repository repositories.Article) *Article {
	return &Article{
		repository: repository,
	}
}

func (instance Article) GetArticles(filter filters.ArticleFilter, userId uuid.UUID) ([]article.Article, int, error) {
	return instance.repository.GetArticles(filter, userId)
}

func (instance Article) GetTrendingArticles(filter filters.ArticleFilter, userId uuid.UUID) ([]article.Article, int, error) {
	return instance.repository.GetTrendingArticles(filter, userId)
}

func (instance Article) GetTrendingArticlesByPropositionType(propositionTypeId uuid.UUID, itemsPerType int,
	userId uuid.UUID) ([]article.Article, error) {
	return instance.repository.GetTrendingArticlesByPropositionType(propositionTypeId, itemsPerType, userId)
}

func (instance Article) GetArticlesToViewLater(filter filters.ArticleFilter, userId uuid.UUID) ([]article.Article, int, error) {
	return instance.repository.GetArticlesToViewLater(filter, userId)
}

func (instance Article) GetPropositionArticlesByNewsletterId(newsletterId uuid.UUID, userId uuid.UUID) ([]article.Article, error) {
	return instance.repository.GetPropositionArticlesByNewsletterId(newsletterId, userId)
}

func (instance Article) GetNewsletterArticleByPropositionId(propositionId uuid.UUID, userId uuid.UUID) (*article.Article, error) {
	return instance.repository.GetNewsletterArticleByPropositionId(propositionId, userId)
}

func (instance Article) SaveArticleRating(userId uuid.UUID, articleId uuid.UUID, rating int) error {
	return instance.repository.SaveArticleRating(userId, articleId, rating)
}

func (instance Article) SaveArticleToViewLater(userId uuid.UUID, articleId uuid.UUID, viewLater bool) error {
	return instance.repository.SaveArticleToViewLater(userId, articleId, viewLater)
}
