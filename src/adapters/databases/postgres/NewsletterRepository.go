package postgres

import (
	"fmt"
	"github.com/devlucassantos/vnc-domains/src/domains/article"
	"github.com/devlucassantos/vnc-domains/src/domains/articlesituation"
	"github.com/devlucassantos/vnc-domains/src/domains/articletype"
	"github.com/devlucassantos/vnc-domains/src/domains/newsletter"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"vnc-api/adapters/databases/dto"
	"vnc-api/adapters/databases/postgres/queries"
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
		log.Errorf("Error retrieving newsletter data of article %s from the database: %s", articleId, err.Error())
		return nil, err
	}

	var newsletterArticles []dto.Article
	err = postgresConnection.Select(&newsletterArticles, queries.Article().Select().ArticlesByNewsletterId(),
		newsletterArticle.Newsletter.Id)
	if err != nil {
		log.Errorf("Error retrieving data for articles of the newsletter %s from the database: %s",
			articleId, err.Error())
		return nil, err
	}

	userArticles := make(map[uuid.UUID]dto.UserArticle)
	if userId != uuid.Nil {
		var articles []interface{}
		articles = append(articles, newsletterArticle.Id)
		for _, articleData := range newsletterArticles {
			articles = append(articles, articleData.Id)
		}
		articleFilters := append([]interface{}{}, userId)
		articleFilters = append(articleFilters, articles...)

		var userArticleData []dto.UserArticle
		err = postgresConnection.Select(&userArticleData, queries.Article().Select().
			RatingsAndArticlesSavedForLaterViewing(len(articles)), articleFilters...)
		if err != nil {
			log.Errorf("Error retrieving articles that the user %s rated and/or saved for later viewing: %s",
				userId, err.Error())
			return nil, err
		}

		for _, userArticle := range userArticleData {
			userArticles[userArticle.Article.Id] = userArticle
		}
	}

	var articles []article.Article
	for _, articleData := range newsletterArticles {
		articleType, err := articletype.NewBuilder().
			Id(articleData.ArticleType.Id).
			Description(articleData.ArticleType.Description).
			Codes(articleData.ArticleType.Codes).
			Color(articleData.ArticleType.Color).
			Build()
		if err != nil {
			log.Errorf("Error validating data for article type %s of article %s: %s",
				articleData.ArticleType.Id, articleData.Id, err.Error())
			return nil, err
		}

		articleBuilder := article.NewBuilder().
			Id(articleData.Id).
			AverageRating(articleData.AverageRating).
			NumberOfRatings(articleData.NumberOfRatings).
			Type(*articleType).
			CreatedAt(articleData.CreatedAt).
			UpdatedAt(articleData.UpdatedAt)

		if _, exists := userArticles[articleData.Id]; exists {
			articleBuilder.UserRating(userArticles[articleData.Id].Rating).
				ViewLater(userArticles[articleData.Id].ViewLater)
		}

		var articleDomain *article.Article
		var articleErr error
		if articleData.Proposition.Id != uuid.Nil {
			articleSpecificType, err := articletype.NewBuilder().
				Id(articleData.Proposition.PropositionType.Id).
				Description(articleData.Proposition.PropositionType.Description).
				Color(articleData.Proposition.PropositionType.Color).
				Build()
			if err != nil {
				log.Errorf("Error validating data for proposition type %s of proposition %s of article %s: %s",
					articleData.Proposition.PropositionType.Id, articleData.Proposition.Id, articleData.Id, err.Error())
				return nil, err
			}

			if articleData.Proposition.ImageUrl != "" {
				articleBuilder.MultimediaUrl(articleData.Proposition.ImageUrl).
					MultimediaDescription(articleData.Proposition.ImageDescription)
			}

			articleDomain, articleErr = articleBuilder.
				Title(articleData.Proposition.Title).
				Content(articleData.Proposition.Content).
				SpecificType(*articleSpecificType).
				Build()
		} else if articleData.Voting.Id != uuid.Nil {
			articleSituation, err := articlesituation.NewBuilder().
				Result(articleData.Voting.Result).
				ResultAnnouncedAt(articleData.Voting.ResultAnnouncedAt).
				IsApproved(articleData.Voting.IsApproved).
				Build()
			if err != nil {
				log.Errorf("Error validating data for article/voting situation of voting %s of article %s: %s",
					articleData.Voting.Id, articleData.Id, err.Error())
				return nil, err
			}

			articleDomain, articleErr = articleBuilder.
				Title(fmt.Sprint("Votação ", articleData.Voting.Code)).
				Content(articleData.Voting.Description).
				Situation(*articleSituation).
				Build()
		} else if articleData.Event.Id != uuid.Nil {
			if articleData.Event.VideoUrl != "" {
				articleBuilder.MultimediaUrl(articleData.Event.VideoUrl)
			}

			articleSituationBuilder := articlesituation.NewBuilder()

			if !articleData.Event.EndsAt.IsZero() {
				articleSituationBuilder.EndsAt(articleData.Event.EndsAt)
			}

			articleSituation, err := articleSituationBuilder.
				Id(articleData.Event.EventSituation.Id).
				Description(articleData.Event.EventSituation.Description).
				Color(articleData.Event.EventSituation.Color).
				StartsAt(articleData.Event.StartsAt).
				Build()
			if err != nil {
				log.Errorf("Error validating data for article/event situation %s of event %s of article %s: %s",
					articleData.Event.EventSituation.Id, articleData.Event.Id, articleData.Id, err.Error())
				return nil, err
			}

			articleSpecificType, err := articletype.NewBuilder().
				Id(articleData.Event.EventType.Id).
				Description(articleData.Event.EventType.Description).
				Color(articleData.Event.EventType.Color).
				Build()
			if err != nil {
				log.Errorf("Error validating data for event type %s of event %s of article %s: %s",
					articleData.Event.EventType.Id, articleData.Event.Id, articleData.Id, err.Error())
				return nil, err
			}

			articleDomain, articleErr = articleBuilder.
				Title(articleData.Event.Title).
				Content(articleData.Event.Description).
				Situation(*articleSituation).
				SpecificType(*articleSpecificType).
				Build()
		} else {
			articleDomain, articleErr = articleBuilder.
				Title(fmt.Sprint("Boletim do dia ",
					articleData.Newsletter.ReferenceDate.Format("02/01/2006"))).
				Content(articleData.Newsletter.Description).
				Build()
		}
		if articleErr != nil {
			log.Errorf("Error validating data for article %s: %s", articleData.Id, articleErr.Error())
			return nil, articleErr
		}

		articles = append(articles, *articleDomain)
	}

	articleType, err := articletype.NewBuilder().
		Id(newsletterArticle.ArticleType.Id).
		Description(newsletterArticle.ArticleType.Description).
		Codes(newsletterArticle.ArticleType.Codes).
		Color(newsletterArticle.ArticleType.Color).
		Build()
	if err != nil {
		log.Errorf("Error validating data for article type %s of article %s: %s",
			newsletterArticle.ArticleType.Id, articleId, err.Error())
		return nil, err
	}

	articleBuilder := article.NewBuilder()

	if _, exists := userArticles[articleId]; exists {
		articleBuilder.UserRating(userArticles[articleId].Rating).ViewLater(userArticles[articleId].ViewLater)
	}

	articleDomain, err := articleBuilder.
		Id(newsletterArticle.Id).
		AverageRating(newsletterArticle.AverageRating).
		NumberOfRatings(newsletterArticle.NumberOfRatings).
		Type(*articleType).
		CreatedAt(newsletterArticle.CreatedAt).
		UpdatedAt(newsletterArticle.UpdatedAt).
		Build()
	if err != nil {
		log.Errorf("Error validating data for article %s of newsletter %s: %s", articleId,
			newsletterArticle.Newsletter.Id, err.Error())
		return nil, err
	}

	newsletterDomain, err := newsletter.NewBuilder().
		Id(newsletterArticle.Newsletter.Id).
		ReferenceDate(newsletterArticle.Newsletter.ReferenceDate).
		Title(fmt.Sprint("Boletim do dia ", newsletterArticle.ReferenceDate.Format("02/01/2006"))).
		Description(newsletterArticle.Newsletter.Description).
		Article(*articleDomain).
		Articles(articles).
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

	_, err = postgresConnection.Exec(queries.ArticleView().Insert(), articleId, userIdPointer)
	if err != nil {
		log.Errorf("Error registering the view for article %s: %s", articleId, err.Error())
	}

	return newsletterDomain, nil
}
