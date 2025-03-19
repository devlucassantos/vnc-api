package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/devlucassantos/vnc-domains/src/domains/article"
	"github.com/devlucassantos/vnc-domains/src/domains/articlesituation"
	"github.com/devlucassantos/vnc-domains/src/domains/articletype"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"vnc-api/core/filters"
	"vnc-api/infra/dto"
	"vnc-api/infra/postgres/queries"
)

type Article struct {
	connectionManager connectionManagerInterface
}

func NewArticleRepository(connectionManager connectionManagerInterface) *Article {
	return &Article{
		connectionManager: connectionManager,
	}
}

func (instance Article) GetArticles(filter filters.Article, userId uuid.UUID) ([]article.Article, int, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("Error creating a connection to the Postgres database: ", err.Error())
		return nil, 0, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var articles []dto.Article
	if !filter.Proposition.IsZero() {
		err = postgresConnection.Select(&articles, queries.Article().Select().Propositions(),
			filter.TypeId, filter.SpecificTypeId, fmt.Sprint("%", filter.Content, "%"), filter.StartDate,
			filter.EndDate, filter.Proposition.DeputyId, filter.Proposition.PartyId,
			filter.Proposition.ExternalAuthorId, filter.Pagination.CalculateOffset(),
			filter.Pagination.GetItemsPerPage())
	} else if !filter.Voting.IsZero() {
		err = postgresConnection.Select(&articles, queries.Article().Select().Votes(),
			filter.TypeId, filter.SpecificTypeId, fmt.Sprint("%", filter.Content, "%"), filter.StartDate,
			filter.EndDate, filter.Voting.StartDate, filter.Voting.EndDate, filter.Voting.IsVotingApproved,
			filter.Voting.LegislativeBodyId, filter.Pagination.CalculateOffset(), filter.Pagination.GetItemsPerPage())
	} else if !filter.Event.IsZero() {
		err = postgresConnection.Select(&articles, queries.Article().Select().Events(),
			filter.TypeId, filter.SpecificTypeId, fmt.Sprint("%", filter.Content, "%"), filter.StartDate,
			filter.EndDate, filter.Event.StartDate, filter.Event.EndDate, filter.Event.SituationId,
			filter.Event.LegislativeBodyId, filter.Event.RapporteurId, filter.Pagination.CalculateOffset(),
			filter.Pagination.GetItemsPerPage())
	} else {
		err = postgresConnection.Select(&articles, queries.Article().Select().All(),
			filter.TypeId, filter.SpecificTypeId, fmt.Sprint("%", filter.Content, "%"), filter.StartDate,
			filter.EndDate, filter.Pagination.CalculateOffset(), filter.Pagination.GetItemsPerPage())
	}
	if err != nil {
		log.Error("Error retrieving data for articles from the database: ", err.Error())
		return nil, 0, err
	}

	userArticles := make(map[uuid.UUID]dto.UserArticle)
	if userId != uuid.Nil && articles != nil {
		var articleFilters []interface{}
		articleFilters = append(articleFilters, userId)
		for _, articleData := range articles {
			articleFilters = append(articleFilters, articleData.Id)
		}

		var userArticleData []dto.UserArticle
		err = postgresConnection.Select(&userArticleData,
			queries.Article().Select().RatingsAndArticlesSavedForLaterViewing(len(articles)), articleFilters...)
		if err != nil {
			log.Errorf("Error retrieving articles that the user %s rated and/or saved for later viewing: %s",
				userId, err.Error())
			return nil, 0, err
		}

		for _, userArticle := range userArticleData {
			userArticles[userArticle.Article.Id] = userArticle
		}
	}

	var articleSlice []article.Article
	for _, articleData := range articles {
		articleType, err := articletype.NewBuilder().
			Id(articleData.ArticleType.Id).
			Description(articleData.ArticleType.Description).
			Codes(articleData.ArticleType.Codes).
			Color(articleData.ArticleType.Color).
			Build()
		if err != nil {
			log.Errorf("Error validating data for article type %s of article %s: %s", articleData.ArticleType.Id,
				articleData.Id, err.Error())
			return nil, 0, err
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
		if articleData.Proposition != nil && articleData.Proposition.Id != uuid.Nil {
			articleSpecificType, err := articletype.NewBuilder().
				Id(articleData.Proposition.PropositionType.Id).
				Description(articleData.Proposition.PropositionType.Description).
				Color(articleData.Proposition.PropositionType.Color).
				Build()
			if err != nil {
				log.Errorf("Error validating data for proposition type %s of proposition %s of article %s: %s",
					articleData.Proposition.PropositionType.Id, articleData.Proposition.Id, articleData.Id, err.Error())
				return nil, 0, err
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
		} else if articleData.Voting != nil && articleData.Voting.Id != uuid.Nil {
			articleSituation, err := articlesituation.NewBuilder().
				IsApproved(articleData.Voting.IsApproved).
				Build()
			if err != nil {
				log.Errorf("Error validating data for article/voting situation of voting %s of article %s: %s",
					articleData.Voting.Id, articleData.Id, err.Error())
				return nil, 0, err
			}

			articleDomain, articleErr = articleBuilder.
				Title(fmt.Sprint("Votação ", articleData.Voting.Code)).
				Content(articleData.Voting.Result).
				Situation(*articleSituation).
				Build()
		} else if articleData.Event != nil && articleData.Event.Id != uuid.Nil {
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
				return nil, 0, err
			}

			articleSpecificType, err := articletype.NewBuilder().
				Id(articleData.Event.EventType.Id).
				Description(articleData.Event.EventType.Description).
				Color(articleData.Event.EventType.Color).
				Build()
			if err != nil {
				log.Errorf("Error validating data for event type %s of event %s of article %s: %s",
					articleData.Event.EventType.Id, articleData.Event.Id, articleData.Id, err.Error())
				return nil, 0, err
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
			return nil, 0, articleErr
		}

		articleSlice = append(articleSlice, *articleDomain)
	}

	var totalNumberOfArticles int
	if !filter.Proposition.IsZero() {
		err = postgresConnection.Get(&totalNumberOfArticles, queries.Article().Select().TotalNumberOfPropositions(),
			filter.TypeId, filter.SpecificTypeId, fmt.Sprint("%", filter.Content, "%"), filter.StartDate,
			filter.EndDate, filter.Proposition.DeputyId, filter.Proposition.PartyId,
			filter.Proposition.ExternalAuthorId)
	} else if !filter.Voting.IsZero() {
		err = postgresConnection.Get(&totalNumberOfArticles, queries.Article().Select().TotalNumberOfVotes(),
			filter.TypeId, filter.SpecificTypeId, fmt.Sprint("%", filter.Content, "%"), filter.StartDate,
			filter.EndDate, filter.Voting.StartDate, filter.Voting.EndDate, filter.Voting.IsVotingApproved,
			filter.Voting.LegislativeBodyId)
	} else if !filter.Event.IsZero() {
		err = postgresConnection.Get(&totalNumberOfArticles, queries.Article().Select().TotalNumberOfEvents(),
			filter.TypeId, filter.SpecificTypeId, fmt.Sprint("%", filter.Content, "%"), filter.StartDate,
			filter.EndDate, filter.Event.StartDate, filter.Event.EndDate, filter.Event.SituationId,
			filter.Event.LegislativeBodyId, filter.Event.RapporteurId)
	} else {
		err = postgresConnection.Get(&totalNumberOfArticles, queries.Article().Select().TotalNumberOfArticles(),
			filter.TypeId, filter.SpecificTypeId, fmt.Sprint("%", filter.Content, "%"), filter.StartDate,
			filter.EndDate)
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Error("Error retrieving the total number of articles from the database: ", err.Error())
		return nil, 0, err
	}

	return articleSlice, totalNumberOfArticles, nil
}

func (instance Article) GetTrendingArticles(filter filters.Article, userId uuid.UUID) ([]article.Article, int, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("Error creating a connection to the Postgres database: ", err.Error())
		return nil, 0, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var trendingArticles []dto.Article
	if !filter.Proposition.IsZero() {
		err = postgresConnection.Select(&trendingArticles, queries.Article().Select().TrendingPropositions(),
			filter.TypeId, filter.SpecificTypeId, fmt.Sprint("%", filter.Content, "%"), filter.StartDate,
			filter.EndDate, filter.Proposition.DeputyId, filter.Proposition.PartyId,
			filter.Proposition.ExternalAuthorId, filter.Pagination.CalculateOffset(),
			filter.Pagination.GetItemsPerPage())
	} else if !filter.Voting.IsZero() {
		err = postgresConnection.Select(&trendingArticles, queries.Article().Select().TrendingVotes(),
			filter.TypeId, filter.SpecificTypeId, fmt.Sprint("%", filter.Content, "%"), filter.StartDate,
			filter.EndDate, filter.Voting.StartDate, filter.Voting.EndDate, filter.Voting.IsVotingApproved,
			filter.Voting.LegislativeBodyId, filter.Pagination.CalculateOffset(), filter.Pagination.GetItemsPerPage())
	} else if !filter.Event.IsZero() {
		err = postgresConnection.Select(&trendingArticles, queries.Article().Select().TrendingEvents(),
			filter.TypeId, filter.SpecificTypeId, fmt.Sprint("%", filter.Content, "%"), filter.StartDate,
			filter.EndDate, filter.Event.StartDate, filter.Event.EndDate, filter.Event.SituationId,
			filter.Event.LegislativeBodyId, filter.Event.RapporteurId, filter.Pagination.CalculateOffset(),
			filter.Pagination.GetItemsPerPage())
	} else {
		err = postgresConnection.Select(&trendingArticles, queries.Article().Select().TrendingArticles(),
			filter.TypeId, filter.SpecificTypeId, fmt.Sprint("%", filter.Content, "%"), filter.StartDate,
			filter.EndDate, filter.Pagination.CalculateOffset(), filter.Pagination.GetItemsPerPage())
	}
	if err != nil {
		log.Error("Error searching for trending articles in the database: ", err.Error())
		return nil, 0, err
	}

	userArticles := make(map[uuid.UUID]dto.UserArticle)
	if userId != uuid.Nil && trendingArticles != nil {
		var articleFilters []interface{}
		articleFilters = append(articleFilters, userId)
		for _, articleData := range trendingArticles {
			articleFilters = append(articleFilters, articleData.Id)
		}

		var userArticleData []dto.UserArticle
		err = postgresConnection.Select(&userArticleData,
			queries.Article().Select().RatingsAndArticlesSavedForLaterViewing(len(trendingArticles)), articleFilters...)
		if err != nil {
			log.Errorf("Error retrieving articles that the user %s rated and/or saved for later viewing: %s",
				userId, err.Error())
			return nil, 0, err
		}

		for _, userArticle := range userArticleData {
			userArticles[userArticle.Article.Id] = userArticle
		}
	}

	var articles []article.Article
	for _, articleData := range trendingArticles {
		articleType, err := articletype.NewBuilder().
			Id(articleData.ArticleType.Id).
			Description(articleData.ArticleType.Description).
			Codes(articleData.ArticleType.Codes).
			Color(articleData.ArticleType.Color).
			Build()
		if err != nil {
			log.Errorf("Error validating data for article type %s of article %s: %s", articleData.ArticleType.Id,
				articleData.Id, err.Error())
			return nil, 0, err
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
		if articleData.Proposition != nil && articleData.Proposition.Id != uuid.Nil {
			articleSpecificType, err := articletype.NewBuilder().
				Id(articleData.Proposition.PropositionType.Id).
				Description(articleData.Proposition.PropositionType.Description).
				Color(articleData.Proposition.PropositionType.Color).
				Build()
			if err != nil {
				log.Errorf("Error validating data for proposition type %s of proposition %s of article %s: %s",
					articleData.Proposition.PropositionType.Id, articleData.Proposition.Id, articleData.Id, err.Error())
				return nil, 0, err
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
		} else if articleData.Voting != nil && articleData.Voting.Id != uuid.Nil {
			articleSituation, err := articlesituation.NewBuilder().
				IsApproved(articleData.Voting.IsApproved).
				Build()
			if err != nil {
				log.Errorf("Error validating data for article/voting situation of voting %s of article %s: %s",
					articleData.Voting.Id, articleData.Id, err.Error())
				return nil, 0, err
			}

			articleDomain, articleErr = articleBuilder.
				Title(fmt.Sprint("Votação ", articleData.Voting.Code)).
				Content(articleData.Voting.Result).
				Situation(*articleSituation).
				Build()
		} else if articleData.Event != nil && articleData.Event.Id != uuid.Nil {
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
				return nil, 0, err
			}

			articleSpecificType, err := articletype.NewBuilder().
				Id(articleData.Event.EventType.Id).
				Description(articleData.Event.EventType.Description).
				Color(articleData.Event.EventType.Color).
				Build()
			if err != nil {
				log.Errorf("Error validating data for event type %s of event %s of article %s: %s",
					articleData.Event.EventType.Id, articleData.Event.Id, articleData.Id, err.Error())
				return nil, 0, err
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
			return nil, 0, articleErr
		}

		articles = append(articles, *articleDomain)
	}

	var totalNumberOfArticles int
	if !filter.Proposition.IsZero() {
		err = postgresConnection.Get(&totalNumberOfArticles, queries.Article().Select().TotalNumberOfPropositions(),
			filter.TypeId, filter.SpecificTypeId, fmt.Sprint("%", filter.Content, "%"), filter.StartDate,
			filter.EndDate, filter.Proposition.DeputyId, filter.Proposition.PartyId, filter.Proposition.ExternalAuthorId)
	} else if !filter.Voting.IsZero() {
		err = postgresConnection.Get(&totalNumberOfArticles, queries.Article().Select().TotalNumberOfVotes(),
			filter.TypeId, filter.SpecificTypeId, fmt.Sprint("%", filter.Content, "%"), filter.StartDate,
			filter.EndDate, filter.Voting.StartDate, filter.Voting.EndDate, filter.Voting.IsVotingApproved,
			filter.Voting.LegislativeBodyId)
	} else if !filter.Event.IsZero() {
		err = postgresConnection.Get(&totalNumberOfArticles, queries.Article().Select().TotalNumberOfEvents(),
			filter.TypeId, filter.SpecificTypeId, fmt.Sprint("%", filter.Content, "%"), filter.StartDate,
			filter.EndDate, filter.Event.StartDate, filter.Event.EndDate, filter.Event.SituationId,
			filter.Event.LegislativeBodyId, filter.Event.RapporteurId)
	} else {
		err = postgresConnection.Get(&totalNumberOfArticles, queries.Article().Select().TotalNumberOfArticles(),
			filter.TypeId, filter.SpecificTypeId, fmt.Sprint("%", filter.Content, "%"), filter.StartDate,
			filter.EndDate)
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Error("Error retrieving the total number of articles from the database: ", err.Error())
		return nil, 0, err
	}

	return articles, totalNumberOfArticles, nil
}

func (instance Article) GetTrendingArticlesByTypeId(articleTypeId uuid.UUID, itemsPerType int, userId uuid.UUID) (
	[]article.Article, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("Error creating a connection to the Postgres database: ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var trendingArticles []dto.Article
	err = postgresConnection.Select(&trendingArticles, queries.Article().Select().TrendingArticlesByTypeId(),
		articleTypeId, itemsPerType)
	if err != nil {
		log.Errorf("Error searching for trending articles of article type %s in the database: %s", articleTypeId,
			err.Error())
		return nil, err
	}

	userArticles := make(map[uuid.UUID]dto.UserArticle)
	if userId != uuid.Nil && trendingArticles != nil {
		var articleFilters []interface{}
		articleFilters = append(articleFilters, userId)
		for _, articleData := range trendingArticles {
			articleFilters = append(articleFilters, articleData.Id)
		}

		var userArticleData []dto.UserArticle
		err = postgresConnection.Select(&userArticleData,
			queries.Article().Select().RatingsAndArticlesSavedForLaterViewing(len(trendingArticles)), articleFilters...)
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
	for _, articleData := range trendingArticles {
		articleType, err := articletype.NewBuilder().
			Id(articleData.ArticleType.Id).
			Description(articleData.ArticleType.Description).
			Codes(articleData.ArticleType.Codes).
			Color(articleData.ArticleType.Color).
			Build()
		if err != nil {
			log.Errorf("Error validating data for article type %s of article %s: %s", articleData.ArticleType.Id,
				articleData.Id, err.Error())
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
		if articleData.Proposition != nil && articleData.Proposition.Id != uuid.Nil {
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
		} else if articleData.Voting != nil && articleData.Voting.Id != uuid.Nil {
			articleSituation, err := articlesituation.NewBuilder().
				IsApproved(articleData.Voting.IsApproved).
				Build()
			if err != nil {
				log.Errorf("Error validating data for article/voting situation of voting %s of article %s: %s",
					articleData.Voting.Id, articleData.Id, err.Error())
				return nil, err
			}

			articleDomain, articleErr = articleBuilder.
				Title(fmt.Sprint("Votação ", articleData.Voting.Code)).
				Content(articleData.Voting.Result).
				Situation(*articleSituation).
				Build()
		} else if articleData.Event != nil && articleData.Event.Id != uuid.Nil {
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

	return articles, nil
}

func (instance Article) GetTrendingArticlesBySpecificTypeId(articleSpecificTypeId uuid.UUID, itemsPerType int,
	userId uuid.UUID) ([]article.Article, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("Error creating a connection to the Postgres database: ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var trendingArticles []dto.Article
	err = postgresConnection.Select(&trendingArticles, queries.Article().Select().TrendingArticlesBySpecificTypeId(),
		articleSpecificTypeId, itemsPerType)
	if err != nil {
		log.Errorf("Error searching for trending articles of article specific type %s in the database: %s",
			articleSpecificTypeId, err.Error())
		return nil, err
	}

	userArticles := make(map[uuid.UUID]dto.UserArticle)
	if userId != uuid.Nil && trendingArticles != nil {
		var articleFilters []interface{}
		articleFilters = append(articleFilters, userId)
		for _, articleData := range trendingArticles {
			articleFilters = append(articleFilters, articleData.Id)
		}

		var userArticleData []dto.UserArticle
		err = postgresConnection.Select(&userArticleData,
			queries.Article().Select().RatingsAndArticlesSavedForLaterViewing(len(trendingArticles)), articleFilters...)
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
	for _, articleData := range trendingArticles {
		articleType, err := articletype.NewBuilder().
			Id(articleData.ArticleType.Id).
			Description(articleData.ArticleType.Description).
			Codes(articleData.ArticleType.Codes).
			Color(articleData.ArticleType.Color).
			Build()
		if err != nil {
			log.Errorf("Error validating data for article type %s of article %s: %s", articleData.ArticleType.Id,
				articleData.Id, err.Error())
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
		if articleData.Proposition != nil && articleData.Proposition.Id != uuid.Nil {
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
		} else {
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
		}
		if articleErr != nil {
			log.Errorf("Error validating data for article %s: %s", articleData.Id, articleErr.Error())
			return nil, articleErr
		}

		articles = append(articles, *articleDomain)
	}

	return articles, nil
}

func (instance Article) GetArticlesToViewLater(filter filters.Article, userId uuid.UUID) ([]article.Article, int, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("Error creating a connection to the Postgres database: ", err.Error())
		return nil, 0, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var userArticles []dto.UserArticle
	if !filter.Proposition.IsZero() {
		err = postgresConnection.Select(&userArticles, queries.Article().Select().PropositionsBookmarkedToViewLater(),
			filter.TypeId, filter.SpecificTypeId, fmt.Sprint("%", filter.Content, "%"), filter.StartDate,
			filter.EndDate, filter.Proposition.DeputyId, filter.Proposition.PartyId,
			filter.Proposition.ExternalAuthorId, userId, filter.Pagination.CalculateOffset(),
			filter.Pagination.GetItemsPerPage())
	} else if !filter.Voting.IsZero() {
		err = postgresConnection.Select(&userArticles, queries.Article().Select().VotesBookmarkedToViewLater(),
			filter.TypeId, filter.SpecificTypeId, fmt.Sprint("%", filter.Content, "%"), filter.StartDate,
			filter.EndDate, filter.Voting.StartDate, filter.Voting.EndDate, filter.Voting.IsVotingApproved,
			filter.Voting.LegislativeBodyId, userId, filter.Pagination.CalculateOffset(),
			filter.Pagination.GetItemsPerPage())
	} else if !filter.Event.IsZero() {
		err = postgresConnection.Select(&userArticles, queries.Article().Select().EventsBookmarkedToViewLater(),
			filter.TypeId, filter.SpecificTypeId, fmt.Sprint("%", filter.Content, "%"), filter.StartDate,
			filter.EndDate, filter.Event.StartDate, filter.Event.EndDate, filter.Event.SituationId,
			filter.Event.LegislativeBodyId, filter.Event.RapporteurId, userId, filter.Pagination.CalculateOffset(),
			filter.Pagination.GetItemsPerPage())
	} else {
		err = postgresConnection.Select(&userArticles, queries.Article().Select().ArticlesBookmarkedToViewLater(),
			filter.TypeId, filter.SpecificTypeId, fmt.Sprint("%", filter.Content, "%"), filter.StartDate,
			filter.EndDate, userId, filter.Pagination.CalculateOffset(), filter.Pagination.GetItemsPerPage())
	}
	if err != nil {
		log.Errorf("Error searching for articles bookmarked for later viewing by user %s in the database: %s",
			userId, err.Error())
		return nil, 0, err
	}

	articles := make(map[uuid.UUID]dto.Article)
	if userArticles != nil {
		var articleIds []interface{}
		for _, articleData := range userArticles {
			articleIds = append(articleIds, articleData.Article.Id)
		}

		var articleDtos []dto.Article
		err = postgresConnection.Select(&articleDtos, queries.Article().Select().In(len(userArticles)), articleIds...)
		if err != nil {
			log.Errorf("Error retrieving data for articles bookmarked for later viewing by user %s from the "+
				"database: %s", userId, err.Error())
			return nil, 0, err
		}

		for _, articleDto := range articleDtos {
			articles[articleDto.Id] = articleDto
		}
	}

	var articleSlice []article.Article
	for _, userArticle := range userArticles {
		articleData := articles[userArticle.Article.Id]

		articleType, err := articletype.NewBuilder().
			Id(articleData.ArticleType.Id).
			Description(articleData.ArticleType.Description).
			Codes(articleData.ArticleType.Codes).
			Color(articleData.ArticleType.Color).
			Build()
		if err != nil {
			log.Errorf("Error validating data for article type %s of article %s: %s", articleData.ArticleType.Id,
				articleData.Id, err.Error())
			return nil, 0, err
		}

		articleBuilder := article.NewBuilder().
			Id(articleData.Id).
			AverageRating(articleData.AverageRating).
			NumberOfRatings(articleData.NumberOfRatings).
			Type(*articleType).
			CreatedAt(articleData.CreatedAt).
			UpdatedAt(articleData.UpdatedAt).
			UserRating(userArticle.Rating).
			ViewLater(userArticle.ViewLater)

		var articleDomain *article.Article
		var articleErr error
		if articleData.Proposition != nil && articleData.Proposition.Id != uuid.Nil {
			articleSpecificType, err := articletype.NewBuilder().
				Id(articleData.Proposition.PropositionType.Id).
				Description(articleData.Proposition.PropositionType.Description).
				Color(articleData.Proposition.PropositionType.Color).
				Build()
			if err != nil {
				log.Errorf("Error validating data for proposition type %s of proposition %s of article %s: %s",
					articleData.Proposition.PropositionType.Id, articleData.Proposition.Id, articleData.Id, err.Error())
				return nil, 0, err
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
				IsApproved(articleData.Voting.IsApproved).
				Build()
			if err != nil {
				log.Errorf("Error validating data for article/voting situation of voting %s of article %s: %s",
					articleData.Voting.Id, articleData.Id, err.Error())
				return nil, 0, err
			}

			articleDomain, articleErr = articleBuilder.
				Title(fmt.Sprint("Votação ", articleData.Voting.Code)).
				Content(articleData.Voting.Result).
				Situation(*articleSituation).
				Build()
		} else if articleData.Event != nil && articleData.Event.Id != uuid.Nil {
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
				return nil, 0, err
			}

			articleSpecificType, err := articletype.NewBuilder().
				Id(articleData.Event.EventType.Id).
				Description(articleData.Event.EventType.Description).
				Color(articleData.Event.EventType.Color).
				Build()
			if err != nil {
				log.Errorf("Error validating data for event type %s of event %s of article %s: %s",
					articleData.Event.EventType.Id, articleData.Event.Id, articleData.Id, err.Error())
				return nil, 0, err
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
			return nil, 0, articleErr
		}

		articleSlice = append(articleSlice, *articleDomain)
	}

	var totalNumberOfArticles int
	if !filter.Proposition.IsZero() {
		err = postgresConnection.Get(&totalNumberOfArticles, queries.Article().Select().
			NumberOfPropositionsBookmarkedToViewLater(), filter.TypeId, filter.SpecificTypeId,
			fmt.Sprint("%", filter.Content, "%"), filter.StartDate, filter.EndDate, filter.Proposition.DeputyId,
			filter.Proposition.PartyId, filter.Proposition.ExternalAuthorId, userId)
	} else if !filter.Voting.IsZero() {
		err = postgresConnection.Get(&totalNumberOfArticles, queries.Article().Select().
			NumberOfVotesBookmarkedToViewLater(), filter.TypeId, filter.SpecificTypeId,
			fmt.Sprint("%", filter.Content, "%"), filter.StartDate, filter.EndDate, filter.Voting.StartDate,
			filter.Voting.EndDate, filter.Voting.IsVotingApproved, filter.Voting.LegislativeBodyId, userId)
	} else if !filter.Event.IsZero() {
		err = postgresConnection.Get(&totalNumberOfArticles, queries.Article().Select().
			NumberOfEventsBookmarkedToViewLater(), filter.TypeId, filter.SpecificTypeId,
			fmt.Sprint("%", filter.Content, "%"), filter.StartDate, filter.EndDate, filter.Event.StartDate,
			filter.Event.EndDate, filter.Event.SituationId, filter.Event.LegislativeBodyId, filter.Event.RapporteurId,
			userId)
	} else {
		err = postgresConnection.Get(&totalNumberOfArticles, queries.Article().Select().
			NumberOfArticlesBookmarkedToViewLater(), filter.TypeId, filter.SpecificTypeId,
			fmt.Sprint("%", filter.Content, "%"), filter.StartDate, filter.EndDate, userId)
	}
	if err != nil {
		log.Errorf("Error retrieving the total number of articles bookmarked for later viewing by user %s from "+
			"the database: %s", userId, err.Error())
		return nil, 0, err
	}

	return articleSlice, totalNumberOfArticles, nil
}

func (instance Article) SaveArticleRating(userId uuid.UUID, articleId uuid.UUID, rating *int) error {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("Error creating a connection to the Postgres database: ", err.Error())
		return err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	sqlResult, err := postgresConnection.Exec(queries.UserArticle().Update().Rating(), rating, userId, articleId)
	if err != nil {
		log.Errorf("Error updating the rating for article %s with user %s: %s", articleId, userId, err.Error())
		return err
	}

	rowsAffected, err := sqlResult.RowsAffected()
	if err == nil && rowsAffected == 0 {
		_, err = postgresConnection.Exec(queries.UserArticle().Insert().Rating(), userId, articleId, rating)
		if err != nil {
			log.Errorf("Error inserting the rating for article %s with user %s: %s",
				articleId, userId, err.Error())
			return err
		}
	} else if err != nil {
		log.Errorf("Error retrieving the number of rows affected by the rating update for article %s with "+
			"user %s: %s", articleId, userId, err.Error())
		return err
	}

	return nil
}

func (instance Article) SaveArticleToViewLater(userId uuid.UUID, articleId uuid.UUID, viewLater bool) error {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("Error creating a connection to the Postgres database: ", err.Error())
		return err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	sqlResult, err := postgresConnection.Exec(queries.UserArticle().Update().ViewLater(), viewLater, userId, articleId)
	if err != nil {
		log.Errorf("Error updating view later bookmark for article %s with user %s: %s",
			articleId, userId, err.Error())
		return err
	}

	rowsAffected, err := sqlResult.RowsAffected()
	if err == nil && rowsAffected == 0 {
		_, err = postgresConnection.Exec(queries.UserArticle().Insert().ViewLater(), userId, articleId, viewLater)
		if err != nil {
			log.Errorf("Error inserting view later bookmark for article %s with user %s: %s",
				articleId, userId, err.Error())
			return err
		}
	} else if err != nil {
		log.Errorf("Error retrieving the number of rows affected by the view later bookmark update for article"+
			" %s with user %s: %s", articleId, userId, err.Error())
		return err
	}

	return nil
}
