package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/devlucassantos/vnc-domains/src/domains/article"
	"github.com/devlucassantos/vnc-domains/src/domains/articletype"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"sort"
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

func (instance Article) GetArticles(filter filters.ArticleFilter, userId uuid.UUID) ([]article.Article, int, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("Error creating a connection to the Postgres database: ", err.Error())
		return nil, 0, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var articles []dto.Article
	if filter.DeputyId != nil || filter.PartyId != nil || filter.ExternalAuthorId != nil {
		err = postgresConnection.Select(&articles, queries.Article().Select().Propositions(),
			&filter.TypeId, fmt.Sprint("%", filter.Content, "%"), filter.DeputyId, filter.PartyId,
			filter.ExternalAuthorId, filter.StartDate, filter.EndDate, filter.PaginationFilter.CalculateOffset(),
			filter.PaginationFilter.GetItemsPerPage())
	} else {
		err = postgresConnection.Select(&articles, queries.Article().Select().All(),
			&filter.TypeId, fmt.Sprint("%", filter.Content, "%"), filter.StartDate, filter.EndDate,
			filter.PaginationFilter.CalculateOffset(), filter.PaginationFilter.GetItemsPerPage())
	}
	if err != nil {
		log.Error("Error fetching data for articles from the database: ", err.Error())
		return nil, 0, err
	}

	var userArticles []dto.UserArticle
	if userId != uuid.Nil && articles != nil {
		var articleFilters []interface{}
		articleFilters = append(articleFilters, userId)
		for _, articleData := range articles {
			articleFilters = append(articleFilters, articleData.Id)
		}

		err = postgresConnection.Select(&userArticles, queries.UserArticle().Select().RatingsAndArticlesSavedForLaterViewing(
			len(articles)), articleFilters...)
		if err != nil {
			log.Errorf("Error fetching articles that the user %s rated and/or saved for later viewing: %s",
				userId, err.Error())
			return nil, 0, err
		}
	}

	var articleSlice []article.Article
	for _, articleData := range articles {
		articleType, err := articletype.NewBuilder().
			Id(articleData.ArticleType.Id).
			Description(articleData.ArticleType.Description).
			Color(articleData.ArticleType.Color).
			SortOrder(articleData.ArticleType.SortOrder).
			CreatedAt(articleData.ArticleType.CreatedAt).
			UpdatedAt(articleData.ArticleType.UpdatedAt).
			Build()
		if err != nil {
			log.Errorf("Error validating data for article type %s: %s", articleData.Id, err.Error())
			continue
		}

		articleBuilder := article.NewBuilder().
			Id(articleData.Id).
			AverageRating(articleData.AverageRating).
			NumberOfRatings(articleData.NumberOfRatings).
			Type(*articleType).
			ReferenceDateTime(articleData.ReferenceDateTime).
			CreatedAt(articleData.CreatedAt).
			UpdatedAt(articleData.UpdatedAt)

		if userArticles != nil {
			for _, userArticle := range userArticles {
				if userArticle.Article.Id == articleData.Id {
					articleBuilder.UserRating(userArticle.Rating).ViewLater(userArticle.ViewLater)
				}
			}
		}

		var articleDomain *article.Article
		if articleData.Proposition.Id != uuid.Nil {
			if articleData.Proposition.ImageUrl != "" {
				articleBuilder.ImageUrl(articleData.Proposition.ImageUrl)
			}

			articleDomain, err = articleBuilder.
				Title(articleData.Proposition.Title).
				Content(articleData.Proposition.Content).
				Build()
		} else {
			articleDomain, err = articleBuilder.
				Title(articleData.Newsletter.Title).
				Content(articleData.Newsletter.Description).
				Build()
		}
		if err != nil {
			log.Errorf("Error validating data for article %s: %s", articleData.Id, err.Error())
			continue
		}

		articleSlice = append(articleSlice, *articleDomain)
	}

	var totalNumberOfArticles int
	if filter.DeputyId != nil || filter.PartyId != nil || filter.ExternalAuthorId != nil {
		err = postgresConnection.Get(&totalNumberOfArticles, queries.Article().Select().TotalNumberOfPropositions(),
			&filter.TypeId, fmt.Sprint("%", filter.Content, "%"), &filter.DeputyId, &filter.PartyId,
			&filter.ExternalAuthorId, filter.StartDate, filter.EndDate)
	} else {
		err = postgresConnection.Get(&totalNumberOfArticles, queries.Article().Select().TotalNumberOfArticles(),
			&filter.TypeId, fmt.Sprint("%", filter.Content, "%"), filter.StartDate, filter.EndDate)
	}
	if err != nil {
		log.Error("Error fetching the total number of articles from the database: ", err.Error())
		return nil, 0, err
	}

	return articleSlice, totalNumberOfArticles, nil
}

func (instance Article) GetTrendingArticles(filter filters.ArticleFilter, userId uuid.UUID) ([]article.Article, int, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("Error creating a connection to the Postgres database: ", err.Error())
		return nil, 0, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var trendingArticles []dto.Article
	if filter.DeputyId != nil || filter.PartyId != nil || filter.ExternalAuthorId != nil {
		err = postgresConnection.Select(&trendingArticles, queries.Article().Select().TrendingPropositions(),
			&filter.TypeId, fmt.Sprint("%", filter.Content, "%"), filter.DeputyId, filter.PartyId,
			filter.ExternalAuthorId, filter.StartDate, filter.EndDate, filter.PaginationFilter.CalculateOffset(),
			filter.PaginationFilter.GetItemsPerPage())
	} else {
		err = postgresConnection.Select(&trendingArticles, queries.Article().Select().TrendingArticles(),
			&filter.TypeId, fmt.Sprint("%", filter.Content, "%"), filter.StartDate, filter.EndDate,
			filter.PaginationFilter.CalculateOffset(), filter.PaginationFilter.GetItemsPerPage())
	}
	if err != nil {
		log.Error("Error searching for trending articles in the database: ", err.Error())
		return nil, 0, err
	}

	var articles []dto.Article
	if trendingArticles != nil {
		var articleIds []interface{}
		for _, articleData := range trendingArticles {
			articleIds = append(articleIds, articleData.Id)
		}

		err = postgresConnection.Select(&articles, queries.Article().Select().In(len(trendingArticles)), articleIds...)
		if err != nil {
			log.Error("Error fetching data for trending articles from the database: ", err.Error())
			return nil, 0, err
		}
	}

	articleOrder := make(map[uuid.UUID]int)
	for index, trendingArticle := range trendingArticles {
		articleOrder[trendingArticle.Id] = index
	}

	sort.Slice(articles, func(currentIndex, comparisonIndex int) bool {
		return articleOrder[articles[currentIndex].Id] < articleOrder[articles[comparisonIndex].Id]
	})

	var userArticles []dto.UserArticle
	if userId != uuid.Nil && articles != nil {
		var articleFilters []interface{}
		articleFilters = append(articleFilters, userId)
		for _, articleData := range articles {
			articleFilters = append(articleFilters, articleData.Id)
		}

		err = postgresConnection.Select(&userArticles, queries.UserArticle().Select().RatingsAndArticlesSavedForLaterViewing(
			len(articles)), articleFilters...)
		if err != nil {
			log.Errorf("Error fetching articles that the user %s rated and/or saved for later viewing: %s",
				userId, err.Error())
			return nil, 0, err
		}
	}

	var articleSlice []article.Article
	for _, articleData := range articles {
		articleType, err := articletype.NewBuilder().
			Id(articleData.ArticleType.Id).
			Description(articleData.ArticleType.Description).
			Color(articleData.ArticleType.Color).
			SortOrder(articleData.ArticleType.SortOrder).
			CreatedAt(articleData.ArticleType.CreatedAt).
			UpdatedAt(articleData.ArticleType.UpdatedAt).
			Build()
		if err != nil {
			log.Errorf("Error validating data for article type %s: %s", articleData.Id, err.Error())
			continue
		}

		articleBuilder := article.NewBuilder().
			Id(articleData.Id).
			AverageRating(articleData.AverageRating).
			NumberOfRatings(articleData.NumberOfRatings).
			Type(*articleType).
			ReferenceDateTime(articleData.ReferenceDateTime).
			CreatedAt(articleData.CreatedAt).
			UpdatedAt(articleData.UpdatedAt)

		if userArticles != nil {
			for _, userArticle := range userArticles {
				if userArticle.Article.Id == articleData.Id {
					articleBuilder.UserRating(userArticle.Rating).ViewLater(userArticle.ViewLater)
				}
			}
		}

		var articleDomain *article.Article
		if articleData.Proposition.Id != uuid.Nil {
			if articleData.Proposition.ImageUrl != "" {
				articleBuilder.ImageUrl(articleData.Proposition.ImageUrl)
			}

			articleDomain, err = articleBuilder.
				Title(articleData.Proposition.Title).
				Content(articleData.Proposition.Content).
				Build()
		} else {
			articleDomain, err = articleBuilder.
				Title(articleData.Newsletter.Title).
				Content(articleData.Newsletter.Description).
				Build()
		}
		if err != nil {
			log.Errorf("Error validating data for article %s: %s", articleData.Id, err.Error())
			continue
		}

		articleSlice = append(articleSlice, *articleDomain)
	}

	var totalNumberOfArticles int
	if filter.DeputyId != nil || filter.PartyId != nil || filter.ExternalAuthorId != nil {
		err = postgresConnection.Get(&totalNumberOfArticles, queries.Article().Select().TotalNumberOfPropositions(),
			&filter.TypeId, fmt.Sprint("%", filter.Content, "%"), &filter.DeputyId, &filter.PartyId,
			&filter.ExternalAuthorId, filter.StartDate, filter.EndDate)
	} else {
		err = postgresConnection.Get(&totalNumberOfArticles, queries.Article().Select().TotalNumberOfArticles(),
			&filter.TypeId, fmt.Sprint("%", filter.Content, "%"), filter.StartDate, filter.EndDate)
	}
	if err != nil {
		log.Error("Error fetching the total number of articles from the database: ", err.Error())
		return nil, 0, err
	}

	return articleSlice, totalNumberOfArticles, nil
}

func (instance Article) GetTrendingArticlesByTypeId(articleTypeId uuid.UUID, numberOfArticles int, userId uuid.UUID) ([]article.Article, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("Error creating a connection to the Postgres database: ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var trendingArticles []dto.Article
	err = postgresConnection.Select(&trendingArticles, queries.Article().Select().TrendingArticlesByType(),
		articleTypeId, numberOfArticles)
	if err != nil {
		log.Errorf("Error searching for trending articles of article type %s in the database: %s",
			articleTypeId, err.Error())
		return nil, err
	}

	var articles []dto.Article
	if trendingArticles != nil {
		var articleIds []interface{}
		for _, articleData := range trendingArticles {
			articleIds = append(articleIds, articleData.Id)
		}

		err = postgresConnection.Select(&articles, queries.Article().Select().In(len(trendingArticles)), articleIds...)
		if err != nil {
			log.Errorf("Error fetching data for trending articles of article type %s from the database: %s",
				articleTypeId, err.Error())
			return nil, err
		}
	}

	articleOrder := make(map[uuid.UUID]int)
	for index, trendingArticle := range trendingArticles {
		articleOrder[trendingArticle.Id] = index
	}

	sort.Slice(articles, func(currentIndex, comparisonIndex int) bool {
		return articleOrder[articles[currentIndex].Id] < articleOrder[articles[comparisonIndex].Id]
	})

	var userArticles []dto.UserArticle
	if userId != uuid.Nil && articles != nil {
		var articleFilters []interface{}
		articleFilters = append(articleFilters, userId)
		for _, articleData := range articles {
			articleFilters = append(articleFilters, articleData.Id)
		}

		err = postgresConnection.Select(&userArticles, queries.UserArticle().Select().RatingsAndArticlesSavedForLaterViewing(
			len(articles)), articleFilters...)
		if err != nil {
			log.Errorf("Error fetching articles that the user %s rated and/or saved for later viewing: %s",
				userId, err.Error())
			return nil, err
		}
	}

	var articleSlice []article.Article
	for _, articleData := range articles {
		articleType, err := articletype.NewBuilder().
			Id(articleData.ArticleType.Id).
			Description(articleData.ArticleType.Description).
			Color(articleData.ArticleType.Color).
			SortOrder(articleData.ArticleType.SortOrder).
			CreatedAt(articleData.ArticleType.CreatedAt).
			UpdatedAt(articleData.ArticleType.UpdatedAt).
			Build()
		if err != nil {
			log.Errorf("Error validating data for article type %s: %s", articleData.Id, err.Error())
			continue
		}

		var articleDomain *article.Article
		articleBuilder := article.NewBuilder().
			Id(articleData.Id).
			AverageRating(articleData.AverageRating).
			NumberOfRatings(articleData.NumberOfRatings).
			Type(*articleType).
			ReferenceDateTime(articleData.ReferenceDateTime).
			CreatedAt(articleData.CreatedAt).
			UpdatedAt(articleData.UpdatedAt)

		if userArticles != nil {
			for _, userArticle := range userArticles {
				if userArticle.Article.Id == articleData.Id {
					articleBuilder.UserRating(userArticle.Rating).ViewLater(userArticle.ViewLater)
				}
			}
		}

		if articleData.Proposition.Id != uuid.Nil {
			if articleData.Proposition.ImageUrl != "" {
				articleBuilder.ImageUrl(articleData.Proposition.ImageUrl)
			}

			articleDomain, err = articleBuilder.
				Title(articleData.Proposition.Title).
				Content(articleData.Proposition.Content).
				Build()
		} else {
			articleDomain, err = articleBuilder.
				Title(articleData.Newsletter.Title).
				Content(articleData.Newsletter.Description).
				Build()
		}
		if err != nil {
			log.Errorf("Error validating data for article %s: %s", articleData.Id, err.Error())
			continue
		}

		articleSlice = append(articleSlice, *articleDomain)
	}

	return articleSlice, nil
}

func (instance Article) GetArticlesToViewLater(filter filters.ArticleFilter, userId uuid.UUID) ([]article.Article, int, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("Error creating a connection to the Postgres database: ", err.Error())
		return nil, 0, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var userArticles []dto.UserArticle
	if filter.DeputyId != nil || filter.PartyId != nil || filter.ExternalAuthorId != nil {
		err = postgresConnection.Select(&userArticles, queries.UserArticle().Select().PropositionsSavedToViewLater(),
			&filter.TypeId, fmt.Sprint("%", filter.Content, "%"), filter.DeputyId, filter.PartyId,
			filter.ExternalAuthorId, filter.StartDate, filter.EndDate, userId, filter.PaginationFilter.CalculateOffset(),
			filter.PaginationFilter.GetItemsPerPage())
	} else {
		err = postgresConnection.Select(&userArticles, queries.UserArticle().Select().ArticlesSavedToViewLater(),
			&filter.TypeId, fmt.Sprint("%", filter.Content, "%"), filter.StartDate, filter.EndDate, userId,
			filter.PaginationFilter.CalculateOffset(), filter.PaginationFilter.GetItemsPerPage())
	}
	if err != nil {
		log.Errorf("Error searching for articles bookmarked for later viewing by user %s in the database: %s",
			userId, err.Error())
		return nil, 0, err
	}

	var articles []dto.Article
	if userArticles != nil {
		var articleIds []interface{}
		for _, articleData := range userArticles {
			articleIds = append(articleIds, articleData.Article.Id)
		}

		err = postgresConnection.Select(&articles, queries.Article().Select().In(len(userArticles)), articleIds...)
		if err != nil {
			log.Errorf("Error fetching data for articles bookmarked for later viewing by user %s from the "+
				"database: %s", userId, err.Error())
			return nil, 0, err
		}
	}

	articleOrder := make(map[uuid.UUID]int)
	for index, userArticle := range userArticles {
		articleOrder[userArticle.Article.Id] = index
	}

	sort.Slice(articles, func(currentIndex, comparisonIndex int) bool {
		return articleOrder[articles[currentIndex].Id] < articleOrder[articles[comparisonIndex].Id]
	})

	var articleSlice []article.Article
	for _, articleData := range articles {
		articleType, err := articletype.NewBuilder().
			Id(articleData.ArticleType.Id).
			Description(articleData.ArticleType.Description).
			Color(articleData.ArticleType.Color).
			SortOrder(articleData.ArticleType.SortOrder).
			CreatedAt(articleData.ArticleType.CreatedAt).
			UpdatedAt(articleData.ArticleType.UpdatedAt).
			Build()
		if err != nil {
			log.Errorf("Error validating data for article type %s: %s", articleData.Id, err.Error())
			continue
		}

		articleBuilder := article.NewBuilder().
			Id(articleData.Id).
			AverageRating(articleData.AverageRating).
			NumberOfRatings(articleData.NumberOfRatings).
			Type(*articleType).
			ReferenceDateTime(articleData.ReferenceDateTime).
			CreatedAt(articleData.CreatedAt).
			UpdatedAt(articleData.UpdatedAt)

		if userArticles != nil {
			for _, userArticle := range userArticles {
				if userArticle.Article.Id == articleData.Id {
					articleBuilder.UserRating(userArticle.Rating).ViewLater(userArticle.ViewLater)
				}
			}
		}

		var articleDomain *article.Article
		if articleData.Proposition.Id != uuid.Nil {
			if articleData.Proposition.ImageUrl != "" {
				articleBuilder.ImageUrl(articleData.Proposition.ImageUrl)
			}

			articleDomain, err = articleBuilder.
				Title(articleData.Proposition.Title).
				Content(articleData.Proposition.Content).
				Build()
		} else {
			articleDomain, err = articleBuilder.
				Title(articleData.Newsletter.Title).
				Content(articleData.Newsletter.Description).
				Build()
		}
		if err != nil {
			log.Errorf("Error validating data for article %s: %s", articleData.Id, err.Error())
			continue
		}

		articleSlice = append(articleSlice, *articleDomain)
	}

	var totalNumberOfArticles int
	if filter.DeputyId != nil || filter.PartyId != nil || filter.ExternalAuthorId != nil {
		err = postgresConnection.Get(&totalNumberOfArticles, queries.UserArticle().Select().NumberOfPropositionsSavedToViewLater(),
			fmt.Sprint("%", filter.Content, "%"), &filter.DeputyId, &filter.PartyId, &filter.ExternalAuthorId,
			&filter.TypeId, filter.StartDate, filter.EndDate, userId)
	} else {
		err = postgresConnection.Get(&totalNumberOfArticles, queries.UserArticle().Select().NumberOfArticlesSavedToViewLater(),
			&filter.TypeId, fmt.Sprint("%", filter.Content, "%"), filter.StartDate, filter.EndDate, userId)
	}
	if err != nil {
		log.Errorf("Error fetching the total number of articles bookmarked for later viewing by user %s from the "+
			"database: %s", userId, err.Error())
		return nil, 0, err
	}

	return articleSlice, totalNumberOfArticles, nil
}

func (instance Article) GetPropositionArticlesByNewsletterId(newsletterId uuid.UUID, userId uuid.UUID) ([]article.Article, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("Error creating a connection to the Postgres database: ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var propositionArticles []dto.Article
	err = postgresConnection.Select(&propositionArticles, queries.Article().Select().PropositionsByNewsletterId(),
		newsletterId)
	if err != nil {
		log.Errorf("Error fetching the data for the articles of the propositions from bulletin %s in the"+
			"database: %s", newsletterId, err.Error())
		return nil, err
	}

	var userArticles []dto.UserArticle
	if userId != uuid.Nil && propositionArticles != nil {
		var articleFilters []interface{}
		articleFilters = append(articleFilters, userId)
		for _, articleData := range propositionArticles {
			articleFilters = append(articleFilters, articleData.Id)
		}

		err = postgresConnection.Select(&userArticles, queries.UserArticle().Select().RatingsAndArticlesSavedForLaterViewing(
			len(propositionArticles)), articleFilters...)
		if err != nil {
			log.Errorf("Error fetching articles that the user %s rated and/or saved for later viewing: %s",
				userId, err.Error())
			return nil, err
		}
	}

	var articles []article.Article
	for _, articleData := range propositionArticles {
		articleBuilder := article.NewBuilder()

		if userArticles != nil {
			for _, userArticle := range userArticles {
				if userArticle.Article.Id == articleData.Id {
					articleBuilder.UserRating(userArticle.Rating).ViewLater(userArticle.ViewLater)
				}
			}
		}

		articleType, err := articletype.NewBuilder().
			Id(articleData.ArticleType.Id).
			Description(articleData.ArticleType.Description).
			Color(articleData.ArticleType.Color).
			SortOrder(articleData.ArticleType.SortOrder).
			CreatedAt(articleData.ArticleType.CreatedAt).
			UpdatedAt(articleData.ArticleType.UpdatedAt).
			Build()
		if err != nil {
			log.Errorf("Error validating data for article type %s: %s", articleData.Id, err.Error())
			continue
		}

		if articleData.Proposition.ImageUrl != "" {
			articleBuilder.ImageUrl(articleData.Proposition.ImageUrl)
		}

		articleDomain, err := articleBuilder.
			Id(articleData.Id).
			Title(articleData.Proposition.Title).
			Content(articleData.Proposition.Content).
			AverageRating(articleData.AverageRating).
			NumberOfRatings(articleData.NumberOfRatings).
			Type(*articleType).
			ReferenceDateTime(articleData.ReferenceDateTime).
			CreatedAt(articleData.CreatedAt).
			UpdatedAt(articleData.UpdatedAt).
			Build()
		if err != nil {
			log.Errorf("Error validating data for article %s of proposition %s: %s", articleData.Id,
				articleData.Proposition.Id, err.Error())
			continue
		}

		articles = append(articles, *articleDomain)
	}

	return articles, nil
}

func (instance Article) GetNewsletterArticleByPropositionId(propositionId uuid.UUID, userId uuid.UUID) (*article.Article, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("Error creating a connection to the Postgres database: ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var newsletterArticle dto.Article
	err = postgresConnection.Get(&newsletterArticle, queries.Article().Select().NewsletterByPropositionId(), propositionId)
	if err != nil {
		errorMessage := "Error fetching data for newsletter by article %s from the database: %s"
		if errors.Is(err, sql.ErrNoRows) {
			log.Infof(errorMessage, propositionId, "Newsletter not found")
		} else {
			log.Errorf(errorMessage, propositionId, err.Error())
		}

		return nil, err
	}

	var userArticle dto.UserArticle
	if userId != uuid.Nil && newsletterArticle.Id != uuid.Nil {
		numberOfArticles := 1
		err = postgresConnection.Get(&userArticle, queries.UserArticle().Select().RatingsAndArticlesSavedForLaterViewing(
			numberOfArticles), userId, newsletterArticle.Id)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			log.Errorf("Error fetching data for article %s that user %s may have rated or saved for later "+
				"viewing: %s", newsletterArticle.Id, userId, err.Error())
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
		log.Errorf("Error validating data for article type %s: %s", newsletterArticle.Id, err.Error())
		return nil, err
	}

	articleDomain, err := articleBuilder.
		Id(newsletterArticle.Id).
		Title(newsletterArticle.Newsletter.Title).
		Content(newsletterArticle.Newsletter.Description).
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

	return articleDomain, nil
}

func (instance Article) SaveArticleRating(userId uuid.UUID, articleId uuid.UUID, rating int) error {
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
		log.Errorf("Error fetching the number of rows affected by the rating update for article %s with "+
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
		log.Errorf("Error fetching the number of rows affected by the view later bookmark update for article"+
			" %s with user %s: %s", articleId, userId, err.Error())
		return err
	}

	return nil
}
