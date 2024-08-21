package postgres

import (
	"database/sql"
	"github.com/devlucassantos/vnc-domains/src/domains/article"
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
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var newsletterArticle dto.Article
	err = postgresConnection.Get(&newsletterArticle, queries.Newsletter().Select().ByArticleId(), articleId)
	if err != nil {
		log.Errorf("Erro ao obter os dados do boletim pela matéria %s no banco de dados: %s", articleId, err.Error())
		return nil, err
	}

	var userArticle dto.UserArticle
	if userId != uuid.Nil && newsletterArticle.Id != uuid.Nil {
		numberOfArticles := 1
		err = postgresConnection.Get(&userArticle, queries.UserArticle().Select().RatingsAndArticlesSavedForLaterViewing(
			numberOfArticles), userId, newsletterArticle.Id)
		if err != nil && err != sql.ErrNoRows {
			log.Errorf("Erro as buscar no banco de dados as informações da matéria %s que o usuário %s avaliou e/ou "+
				"salvou para ver depois: %s", articleId, userId, err.Error())
			return nil, err
		}
	}

	articleBuilder := article.NewBuilder()

	if userArticle.Article.Id != uuid.Nil {
		articleBuilder.UserRating(userArticle.Rating).ViewLater(userArticle.ViewLater)
	}

	articleDomain, err := articleBuilder.
		Id(newsletterArticle.Id).
		AverageRating(newsletterArticle.AverageRating).
		NumberOfRatings(newsletterArticle.NumberOfRatings).
		Type("Proposição").
		ReferenceDateTime(newsletterArticle.ReferenceDateTime).
		CreatedAt(newsletterArticle.CreatedAt).
		UpdatedAt(newsletterArticle.UpdatedAt).
		Build()
	if err != nil {
		log.Errorf("Erro ao validar os dados da matéria %s do boletim %s: %s", newsletterArticle.Id,
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
		log.Errorf("Erro construindo a estrutura de dados do boletim %s: %s", articleId, err.Error())
		return nil, err
	}

	_, err = postgresConnection.Exec(queries.ArticleView().Insert(), newsletterArticle.Id)
	if err != nil {
		log.Errorf("Erro ao registrar a visualização do boletim %s: %s", articleId, err.Error())
	}

	return newsletterDomain, nil
}
