package postgres

import (
	"fmt"
	"github.com/devlucassantos/vnc-domains/src/domains/article"
	"github.com/devlucassantos/vnc-domains/src/domains/articlesituation"
	"github.com/devlucassantos/vnc-domains/src/domains/articletype"
	"github.com/devlucassantos/vnc-domains/src/domains/deputy"
	"github.com/devlucassantos/vnc-domains/src/domains/externalauthor"
	"github.com/devlucassantos/vnc-domains/src/domains/externalauthortype"
	"github.com/devlucassantos/vnc-domains/src/domains/party"
	"github.com/devlucassantos/vnc-domains/src/domains/proposition"
	"github.com/devlucassantos/vnc-domains/src/domains/propositiontype"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"vnc-api/adapters/databases/dto"
	"vnc-api/adapters/databases/postgres/queries"
)

type Proposition struct {
	connectionManager connectionManagerInterface
}

func NewPropositionRepository(connectionManager connectionManagerInterface) *Proposition {
	return &Proposition{
		connectionManager: connectionManager,
	}
}

func (instance Proposition) GetPropositionByArticleId(articleId uuid.UUID, userId uuid.UUID) (*proposition.Proposition, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("Error creating a connection to the Postgres database: ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var propositionArticle dto.Article
	err = postgresConnection.Get(&propositionArticle, queries.Proposition().Select().ByArticleId(), articleId)
	if err != nil {
		log.Errorf("Error retrieving proposition data for article %s from the database: %s", articleId, err.Error())
		return nil, err
	}

	var propositionAuthors []dto.Author
	err = postgresConnection.Select(&propositionAuthors, queries.PropositionAuthor().Select().ByPropositionId(),
		propositionArticle.Proposition.Id)
	if err != nil {
		log.Errorf("Error retrieving data of the authors for proposition %s from the database: %s", articleId,
			err.Error())
		return nil, err
	}

	var deputies []deputy.Deputy
	var externalAuthors []externalauthor.ExternalAuthor
	for _, author := range propositionAuthors {
		if author.Deputy.Id != uuid.Nil {
			currentParty, err := party.NewBuilder().
				Id(author.Deputy.Party.Id).
				Name(author.Deputy.Party.Name).
				Acronym(author.Deputy.Party.Acronym).
				ImageUrl(author.Deputy.Party.ImageUrl).
				ImageDescription(fmt.Sprintf("Logo do %s (%s)", author.Deputy.Party.Name,
					author.Deputy.Party.Acronym)).
				Build()
			if err != nil {
				log.Errorf("Error validating data for the current party %s of deputy %s for proposition %s: %s",
					author.Deputy.Party.Id, author.Deputy.Id, articleId, err.Error())
			}

			previousParty, err := party.NewBuilder().
				Id(author.Deputy.PreviousParty.Id).
				Name(author.Deputy.PreviousParty.Name).
				Acronym(author.Deputy.PreviousParty.Acronym).
				ImageUrl(author.Deputy.PreviousParty.ImageUrl).
				ImageDescription(fmt.Sprintf("Logo do %s (%s)", author.Deputy.PreviousParty.Name,
					author.Deputy.PreviousParty.Acronym)).
				Build()
			if err != nil {
				log.Errorf("Error validating data for party %s of deputy %s when drafting proposition %s: %s",
					author.Deputy.PreviousParty.Id, author.Deputy.Id, articleId, err.Error())
			}

			deputyDomain, err := deputy.NewBuilder().
				Id(author.Deputy.Id).
				Name(author.Deputy.Name).
				ElectoralName(author.Deputy.ElectoralName).
				ImageUrl(author.Deputy.ImageUrl).
				ImageDescription(fmt.Sprintf("Foto do(a) deputado(a) federal %s (%s-%s)", author.Deputy.Name,
					author.Deputy.Party.Acronym, author.Deputy.FederatedUnit)).
				Party(*currentParty).
				FederatedUnit(author.Deputy.FederatedUnit).
				PreviousParty(*previousParty).
				PreviousFederatedUnit(author.Deputy.PreviousFederatedUnit).
				Build()
			if err != nil {
				log.Errorf("Error validating data for deputy %s for proposition %s: %s", author.Deputy.Id,
					articleId, err.Error())
				return nil, err
			}

			deputies = append(deputies, *deputyDomain)
		} else if author.ExternalAuthor.Id != uuid.Nil {
			externalAuthorType, err := externalauthortype.NewBuilder().
				Id(author.ExternalAuthor.ExternalAuthorType.Id).
				Description(author.ExternalAuthor.ExternalAuthorType.Description).
				Build()
			if err != nil {
				log.Errorf("Error validating data for external author type %s for external author %s: %s",
					author.ExternalAuthor.ExternalAuthorType.Id, author.ExternalAuthor.Id, err.Error())
				return nil, err
			}

			externalAuthor, err := externalauthor.NewBuilder().
				Id(author.ExternalAuthor.Id).
				Name(author.ExternalAuthor.Name).
				Type(*externalAuthorType).
				Build()
			if err != nil {
				log.Errorf("Error validating data for external author %s for proposition %s: %s",
					author.ExternalAuthor.Id, articleId, err.Error())
				return nil, err
			}

			externalAuthors = append(externalAuthors, *externalAuthor)
		}
	}

	var relatedArticleIds []uuid.UUID
	err = postgresConnection.Select(&relatedArticleIds, queries.Article().Select().RelatedArticlesByPropositionId(),
		propositionArticle.Proposition.Id)
	if err != nil {
		log.Errorf("Error retrieving articles IDs related to proposition %s from the database: %s", articleId,
			err.Error())
		return nil, err
	}

	var relatedArticles []dto.Article
	if relatedArticleIds != nil {
		var relatedArticleIdsAsInterfaceSlice []interface{}
		for _, relatedArticleId := range relatedArticleIds {
			relatedArticleIdsAsInterfaceSlice = append(relatedArticleIdsAsInterfaceSlice, relatedArticleId)
		}

		err = postgresConnection.Select(&relatedArticles, queries.Article().Select().In(len(relatedArticleIds)),
			relatedArticleIdsAsInterfaceSlice...)
		if err != nil {
			log.Errorf("Error retrieving data for articles related to proposition %s from the database: %s",
				articleId, err.Error())
			return nil, err
		}
	}

	userArticles := make(map[uuid.UUID]dto.UserArticle)
	if userId != uuid.Nil {
		articles := append([]interface{}{}, propositionArticle.Id)
		for _, relatedArticleId := range relatedArticleIds {
			articles = append(articles, relatedArticleId)
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
	for _, articleData := range relatedArticles {
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
		if articleData.Voting.Id != uuid.Nil {
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

	articleBuilder := article.NewBuilder()

	if _, exists := userArticles[articleId]; exists {
		articleBuilder.UserRating(userArticles[articleId].Rating).ViewLater(userArticles[articleId].ViewLater)
	}

	articleType, err := articletype.NewBuilder().
		Id(propositionArticle.ArticleType.Id).
		Description(propositionArticle.ArticleType.Description).
		Codes(propositionArticle.ArticleType.Codes).
		Color(propositionArticle.ArticleType.Color).
		Build()
	if err != nil {
		log.Errorf("Error validating data for article type %s of article %s: %s",
			propositionArticle.ArticleType.Id, articleId, err.Error())
		return nil, err
	}

	articleDomain, err := articleBuilder.
		Id(propositionArticle.Id).
		AverageRating(propositionArticle.AverageRating).
		NumberOfRatings(propositionArticle.NumberOfRatings).
		Type(*articleType).
		CreatedAt(propositionArticle.CreatedAt).
		UpdatedAt(propositionArticle.UpdatedAt).
		Build()
	if err != nil {
		log.Errorf("Error validating data for article %s of proposition %s: %s", articleId,
			propositionArticle.Proposition.Id, err.Error())
		return nil, err
	}

	propositionType, err := propositiontype.NewBuilder().
		Id(propositionArticle.Proposition.PropositionType.Id).
		Description(propositionArticle.Proposition.PropositionType.Description).
		Color(propositionArticle.Proposition.PropositionType.Color).
		Build()
	if err != nil {
		log.Errorf("Error validating data for proposition type %s of proposition %s of article %s: %s",
			propositionArticle.Proposition.PropositionType.Id, propositionArticle.Proposition.Id,
			propositionArticle.Id, err.Error())
		return nil, err
	}

	propositionBuilder := proposition.NewBuilder()

	if propositionArticle.Proposition.ImageUrl != "" {
		propositionBuilder.ImageUrl(propositionArticle.Proposition.ImageUrl).
			ImageDescription(propositionArticle.Proposition.ImageDescription)
	}

	propositionDomain, err := propositionBuilder.
		Id(propositionArticle.Proposition.Id).
		OriginalTextUrl(propositionArticle.Proposition.OriginalTextUrl).
		OriginalTextMimeType(propositionArticle.Proposition.OriginalTextMimeType).
		Title(propositionArticle.Proposition.Title).
		Content(propositionArticle.Proposition.Content).
		SubmittedAt(propositionArticle.Proposition.SubmittedAt).
		Type(*propositionType).
		Deputies(deputies).
		ExternalAuthors(externalAuthors).
		Article(*articleDomain).
		RelatedArticles(articles).
		Build()
	if err != nil {
		log.Errorf("Error validating data for proposition %s of article %s: %s",
			propositionArticle.Proposition.Id, articleId, err.Error())
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

	return propositionDomain, nil
}
