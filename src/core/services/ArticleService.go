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

func (instance Article) GetArticles(filter filters.Article, userId uuid.UUID) ([]article.Article, int, error) {
	return instance.repository.GetArticles(filter, userId)
}

func (instance Article) GetTrendingArticles(filter filters.Article, userId uuid.UUID) ([]article.Article, int, error) {
	return instance.repository.GetTrendingArticles(filter, userId)
}

func (instance Article) GetTrendingArticlesByTypeId(articleTypeId uuid.UUID, itemsPerType int, userId uuid.UUID) (
	[]article.Article, error) {
	return instance.repository.GetTrendingArticlesByTypeId(articleTypeId, itemsPerType, userId)
}

func (instance Article) GetTrendingArticlesBySpecificTypeId(articleSpecificTypeId uuid.UUID, itemsPerType int,
	userId uuid.UUID) ([]article.Article, error) {
	return instance.repository.GetTrendingArticlesBySpecificTypeId(articleSpecificTypeId, itemsPerType, userId)
}

func (instance Article) GetArticlesToViewLater(filter filters.Article, userId uuid.UUID) ([]article.Article, int, error) {
	return instance.repository.GetArticlesToViewLater(filter, userId)
}

func (instance Article) SaveArticleRating(userId uuid.UUID, articleId uuid.UUID, rating *int) error {
	return instance.repository.SaveArticleRating(userId, articleId, rating)
}

func (instance Article) SaveArticleToViewLater(userId uuid.UUID, articleId uuid.UUID, viewLater bool) error {
	return instance.repository.SaveArticleToViewLater(userId, articleId, viewLater)
}
