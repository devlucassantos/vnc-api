package postgres

import (
	"database/sql"
	"errors"
	"github.com/devlucassantos/vnc-domains/src/domains/article"
	"github.com/devlucassantos/vnc-domains/src/domains/articletype"
	"github.com/devlucassantos/vnc-domains/src/domains/newsletter"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"vnc-api/infra/dto"
	"vnc-api/infra/postgres/queries"
)

type Newsletter struct {
	connectionManager connectionManagerInterface
}

func NewNewsletterRepository(connectionManager connectionManagerInterface) *Newsletter {
	return &Newsletter{
		connectionManager: connectionManager,
	}
}

func (instance Newsletter) GetNewsletterByArticleId(articleId uuid.UUID, userId uuid.UUID) (*newsletter.Newsletter, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("Error creating a connection to the Postgres database: ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var newsletterArticle dto.Article
	err = postgresConnection.Get(&newsletterArticle, queries.Newsletter().Select().ByArticleId(), articleId)
	if err != nil {
		log.Errorf("Error fetching newsletter data for article %s from the database: %s", articleId, err.Error())
		return nil, err
	}

	var userArticle dto.UserArticle
	if userId != uuid.Nil && newsletterArticle.Id != uuid.Nil {
		numberOfArticles := 1
		err = postgresConnection.Get(&userArticle, queries.UserArticle().Select().RatingsAndArticlesSavedForLaterViewing(
			numberOfArticles), userId, newsletterArticle.Id)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			log.Errorf("Error fetching data for article %s that user %s may have rated or saved for later "+
				"viewing: %s", articleId, userId, err.Error())
			return nil, err
		}
	}

	articleBuilder := article.NewBuilder()

	if userArticle.Article != nil {
		articleBuilder.UserRating(userArticle.Rating).ViewLater(userArticle.ViewLater)
	}

	articleType, err := articletype.NewBuilder().
		Id(newsletterArticle.ArticleType.Id).
		Description(newsletterArticle.ArticleType.Description).
		Color(newsletterArticle.ArticleType.Color).
		SortOrder(newsletterArticle.ArticleType.SortOrder).
		CreatedAt(newsletterArticle.ArticleType.CreatedAt).
		UpdatedAt(newsletterArticle.ArticleType.UpdatedAt).
		Build()
	if err != nil {
		log.Errorf("Error validating data for article type %s: %s", articleId, err.Error())
		return nil, err
	}

	articleDomain, err := articleBuilder.
		Id(newsletterArticle.Id).
		AverageRating(newsletterArticle.AverageRating).
		NumberOfRatings(newsletterArticle.NumberOfRatings).
		Type(*articleType).
		ReferenceDateTime(newsletterArticle.ReferenceDateTime).
		CreatedAt(newsletterArticle.CreatedAt).
		UpdatedAt(newsletterArticle.UpdatedAt).
		Build()
	if err != nil {
		log.Errorf("Error validating data for article %s of newsletter %s: %s", newsletterArticle.Id,
			newsletterArticle.Newsletter.Id, err.Error())
		return nil, err
	}

	newsletterDomain, err := newsletter.NewBuilder().
		Id(newsletterArticle.Newsletter.Id).
		ReferenceDate(newsletterArticle.Newsletter.ReferenceDate).
		Title(newsletterArticle.Newsletter.Title).
		Description(newsletterArticle.Newsletter.Description).
		Article(*articleDomain).
		CreatedAt(newsletterArticle.Newsletter.CreatedAt).
		UpdatedAt(newsletterArticle.Newsletter.UpdatedAt).
		Build()
	if err != nil {
		log.Errorf("Error validating data for newsletter %s of article %s: %s",
			newsletterArticle.Newsletter.Id, articleId, err.Error())
		return nil, err
	}

	var userIdPointer *uuid.UUID
	if userId != uuid.Nil {
		userIdPointer = &userId
	}

	_, err = postgresConnection.Exec(queries.ArticleView().Insert(), newsletterArticle.Id, userIdPointer)
	if err != nil {
		log.Errorf("Error registering the view for article %s: %s", articleId, err.Error())
	}

	return newsletterDomain, nil
}
