package postgres

import (
	"database/sql"
	"fmt"
	"github.com/devlucassantos/vnc-domains/src/domains/article"
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

func (instance Article) GetArticles(filter filters.ArticleFilter, userId uuid.UUID) ([]article.Article, int, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("Erro ao tentar se conectar com o Postgres: ", err.Error())
		return nil, 0, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var articles []dto.Article
	if filter.Type == "Proposição" || filter.DeputyId != nil || filter.PartyId != nil || filter.ExternalAuthorId != nil {
		err = postgresConnection.Select(&articles, queries.Article().Select().Propositions(),
			fmt.Sprint("%", filter.Content, "%"), filter.DeputyId, filter.PartyId, filter.ExternalAuthorId,
			filter.StartDate, filter.EndDate, filter.PaginationFilter.CalculateOffset(), filter.PaginationFilter.GetItemsPerPage())
	} else if filter.Type == "Boletim" {
		err = postgresConnection.Select(&articles, queries.Article().Select().Newsletters(),
			fmt.Sprint("%", filter.Content, "%"), filter.StartDate, filter.EndDate,
			filter.PaginationFilter.CalculateOffset(), filter.PaginationFilter.GetItemsPerPage())
	} else {
		err = postgresConnection.Select(&articles, queries.Article().Select().All(),
			fmt.Sprint("%", filter.Content, "%"), filter.StartDate, filter.EndDate,
			filter.PaginationFilter.CalculateOffset(), filter.PaginationFilter.GetItemsPerPage())
	}
	if err != nil {
		log.Error("Erro ao obter os dados das matérias no banco de dados: ", err.Error())
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
			log.Error("Erro as buscar no banco de dados as matérias que o usuário avaliou e/ou salvou para ver depois: ", err.Error())
			return nil, 0, err
		}
	}

	var articleList []article.Article
	for _, articleData := range articles {
		var articleDomain *article.Article
		if articleData.Proposition.Id != uuid.Nil {
			articleBuilder := article.NewBuilder()

			if userArticles != nil {
				for _, userArticle := range userArticles {
					if userArticle.Article.Id == articleData.Id {
						articleBuilder.UserRating(userArticle.Rating).ViewLater(userArticle.ViewLater)
					}
				}
			}

			articleDomain, err = articleBuilder.
				Id(articleData.Id).
				Title(articleData.Proposition.Title).
				Content(articleData.Proposition.Content).
				AverageRating(articleData.AverageRating).
				NumberOfRatings(articleData.NumberOfRatings).
				Type("Proposição").
				ReferenceDateTime(articleData.ReferenceDateTime).
				CreatedAt(articleData.CreatedAt).
				UpdatedAt(articleData.UpdatedAt).
				Build()
			if err != nil {
				log.Errorf("Erro ao validar os dados da matéria %s da proposição %s: %s", articleData.Id,
					articleData.Proposition.Id, err.Error())
				continue
			}
		} else {
			articleBuilder := article.NewBuilder()

			if userArticles != nil {
				for _, userArticle := range userArticles {
					if userArticle.Article.Id == articleData.Id {
						articleBuilder.UserRating(userArticle.Rating).ViewLater(userArticle.ViewLater)
					}
				}
			}

			articleDomain, err = articleBuilder.
				Id(articleData.Id).
				Title(articleData.Newsletter.Title).
				Content(articleData.Newsletter.Description).
				AverageRating(articleData.AverageRating).
				NumberOfRatings(articleData.NumberOfRatings).
				Type("Boletim").
				ReferenceDateTime(articleData.ReferenceDateTime).
				CreatedAt(articleData.CreatedAt).
				UpdatedAt(articleData.UpdatedAt).
				Build()
			if err != nil {
				log.Errorf("Erro ao validar os dados da matéria %s do boletim %s: %s", articleData.Id,
					articleData.Newsletter.Id, err.Error())
				continue
			}
		}

		articleList = append(articleList, *articleDomain)
	}

	var totalNumberOfArticles int
	if filter.Type == "Proposição" || filter.DeputyId != nil || filter.PartyId != nil || filter.ExternalAuthorId != nil {
		err = postgresConnection.Get(&totalNumberOfArticles, queries.Article().Select().TotalNumberOfPropositions(),
			fmt.Sprint("%", filter.Content, "%"), &filter.DeputyId, &filter.PartyId, &filter.ExternalAuthorId,
			filter.StartDate, filter.EndDate)
	} else if filter.Type == "Boletim" {
		err = postgresConnection.Get(&totalNumberOfArticles, queries.Article().Select().TotalNumberOfNewsletters(),
			fmt.Sprint("%", filter.Content, "%"), filter.StartDate, filter.EndDate)
	} else {
		err = postgresConnection.Get(&totalNumberOfArticles, queries.Article().Select().TotalNumberOfArticles(),
			fmt.Sprint("%", filter.Content, "%"), filter.StartDate, filter.EndDate)
	}
	if err != nil {
		log.Error("Erro ao obter a quantidade total de matérias no banco de dados: ", err.Error())
		return nil, 0, err
	}

	return articleList, totalNumberOfArticles, nil
}

func (instance Article) GetTrendingArticles(filter filters.ArticleFilter, userId uuid.UUID) ([]article.Article, int, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("Erro ao tentar se conectar com o Postgres: ", err.Error())
		return nil, 0, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var articles []dto.Article
	if filter.Type == "Proposição" || filter.DeputyId != nil || filter.PartyId != nil || filter.ExternalAuthorId != nil {
		err = postgresConnection.Select(&articles, queries.Article().Select().TrendingPropositions(),
			fmt.Sprint("%", filter.Content, "%"), filter.DeputyId, filter.PartyId, filter.ExternalAuthorId,
			filter.StartDate, filter.EndDate, filter.PaginationFilter.CalculateOffset(), filter.PaginationFilter.GetItemsPerPage())
	} else if filter.Type == "Boletim" {
		err = postgresConnection.Select(&articles, queries.Article().Select().TrendingNewsletters(),
			fmt.Sprint("%", filter.Content, "%"), filter.StartDate, filter.EndDate,
			filter.PaginationFilter.CalculateOffset(), filter.PaginationFilter.GetItemsPerPage())
	} else {
		err = postgresConnection.Select(&articles, queries.Article().Select().TrendingArticles(),
			fmt.Sprint("%", filter.Content, "%"), filter.StartDate, filter.EndDate,
			filter.PaginationFilter.CalculateOffset(), filter.PaginationFilter.GetItemsPerPage())
	}
	if err != nil {
		log.Error("Erro ao obter os dados das matérias em alta no banco de dados: ", err.Error())
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
			log.Error("Erro as buscar no banco de dados as matérias que o usuário avaliou e/ou salvou para ver depois: ", err.Error())
			return nil, 0, err
		}
	}

	var articleList []article.Article
	for _, articleData := range articles {
		var articleDomain *article.Article
		if articleData.Proposition.Id != uuid.Nil {
			articleBuilder := article.NewBuilder()

			if userArticles != nil {
				for _, userArticle := range userArticles {
					if userArticle.Article.Id == articleData.Id {
						articleBuilder.UserRating(userArticle.Rating).ViewLater(userArticle.ViewLater)
					}
				}
			}

			articleDomain, err = articleBuilder.
				Id(articleData.Id).
				Title(articleData.Proposition.Title).
				Content(articleData.Proposition.Content).
				AverageRating(articleData.AverageRating).
				NumberOfRatings(articleData.NumberOfRatings).
				Type("Proposição").
				ReferenceDateTime(articleData.ReferenceDateTime).
				CreatedAt(articleData.CreatedAt).
				UpdatedAt(articleData.UpdatedAt).
				Build()
			if err != nil {
				log.Errorf("Erro ao validar os dados da matéria %s da proposição %s: %s", articleData.Id,
					articleData.Proposition.Id, err.Error())
				continue
			}
		} else {
			articleBuilder := article.NewBuilder()

			if userArticles != nil {
				for _, userArticle := range userArticles {
					if userArticle.Article.Id == articleData.Id {
						articleBuilder.UserRating(userArticle.Rating).ViewLater(userArticle.ViewLater)
					}
				}
			}

			articleDomain, err = articleBuilder.
				Id(articleData.Id).
				Title(articleData.Newsletter.Title).
				Content(articleData.Newsletter.Description).
				AverageRating(articleData.AverageRating).
				NumberOfRatings(articleData.NumberOfRatings).
				Type("Boletim").
				ReferenceDateTime(articleData.ReferenceDateTime).
				CreatedAt(articleData.CreatedAt).
				UpdatedAt(articleData.UpdatedAt).
				Build()
			if err != nil {
				log.Errorf("Erro ao validar os dados da matéria %s do boletim %s: %s", articleData.Id,
					articleData.Newsletter.Id, err.Error())
				continue
			}
		}

		articleList = append(articleList, *articleDomain)
	}

	var totalNumberOfArticles int
	if filter.Type == "Proposição" || filter.DeputyId != nil || filter.PartyId != nil || filter.ExternalAuthorId != nil {
		err = postgresConnection.Get(&totalNumberOfArticles, queries.Article().Select().TotalNumberOfPropositions(),
			fmt.Sprint("%", filter.Content, "%"), &filter.DeputyId, &filter.PartyId, &filter.ExternalAuthorId,
			filter.StartDate, filter.EndDate)
	} else if filter.Type == "Boletim" {
		err = postgresConnection.Get(&totalNumberOfArticles, queries.Article().Select().TotalNumberOfNewsletters(),
			fmt.Sprint("%", filter.Content, "%"), filter.StartDate, filter.EndDate)
	} else {
		err = postgresConnection.Get(&totalNumberOfArticles, queries.Article().Select().TotalNumberOfArticles(),
			fmt.Sprint("%", filter.Content, "%"), filter.StartDate, filter.EndDate)
	}
	if err != nil {
		log.Error("Erro ao obter a quantidade total de matérias no banco de dados: ", err.Error())
		return nil, 0, err
	}

	return articleList, totalNumberOfArticles, nil
}

func (instance Article) GetTrendingArticlesByPropositionType(propositionTypeId uuid.UUID, numberOfArticles int, userId uuid.UUID) ([]article.Article, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("Erro ao tentar se conectar com o Postgres: ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var articles []dto.Article
	err = postgresConnection.Select(&articles, queries.Article().Select().TrendingArticlesByPropositionType(),
		propositionTypeId, numberOfArticles)
	if err != nil {
		log.Errorf("Erro ao obter os dados das matérias do tipo de preposição %s no banco de dados: %s",
			propositionTypeId, err.Error())
		return nil, err
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
			log.Error("Erro as buscar no banco de dados as matérias que o usuário avaliou e/ou salvou para ver depois: ", err.Error())
			return nil, err
		}
	}

	var articleList []article.Article
	for _, articleData := range articles {
		articleBuilder := article.NewBuilder()

		if userArticles != nil {
			for _, userArticle := range userArticles {
				if userArticle.Article.Id == articleData.Id {
					articleBuilder.UserRating(userArticle.Rating).ViewLater(userArticle.ViewLater)
				}
			}
		}

		articleDomain, err := articleBuilder.
			Id(articleData.Id).
			Title(articleData.Proposition.Title).
			Content(articleData.Proposition.Content).
			AverageRating(articleData.AverageRating).
			NumberOfRatings(articleData.NumberOfRatings).
			Type("Proposição").
			ReferenceDateTime(articleData.ReferenceDateTime).
			CreatedAt(articleData.CreatedAt).
			UpdatedAt(articleData.UpdatedAt).
			Build()
		if err != nil {
			log.Errorf("Erro ao validar os dados da matéria %s da proposição %s: %s", articleData.Id,
				articleData.Proposition.Id, err.Error())
			continue
		}

		articleList = append(articleList, *articleDomain)
	}

	return articleList, nil
}

func (instance Article) GetArticlesToViewLater(filter filters.ArticleFilter, userId uuid.UUID) ([]article.Article, int, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("Erro ao tentar se conectar com o Postgres: ", err.Error())
		return nil, 0, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var userArticles []dto.UserArticle
	if filter.Type == "Proposição" || filter.DeputyId != nil || filter.PartyId != nil || filter.ExternalAuthorId != nil {
		err = postgresConnection.Select(&userArticles, queries.UserArticle().Select().PropositionsSavedToViewLater(),
			fmt.Sprint("%", filter.Content, "%"), filter.DeputyId, filter.PartyId, filter.ExternalAuthorId,
			filter.StartDate, filter.EndDate, userId, filter.PaginationFilter.CalculateOffset(),
			filter.PaginationFilter.GetItemsPerPage())
	} else if filter.Type == "Boletim" {
		err = postgresConnection.Select(&userArticles, queries.UserArticle().Select().NewslettersSavedToViewLater(),
			fmt.Sprint("%", filter.Content, "%"), filter.StartDate, filter.EndDate, userId,
			filter.PaginationFilter.CalculateOffset(), filter.PaginationFilter.GetItemsPerPage())
	} else {
		err = postgresConnection.Select(&userArticles, queries.UserArticle().Select().ArticlesSavedToViewLater(),
			fmt.Sprint("%", filter.Content, "%"), filter.StartDate, filter.EndDate, userId,
			filter.PaginationFilter.CalculateOffset(), filter.PaginationFilter.GetItemsPerPage())
	}
	if err != nil {
		log.Error("Erro ao obter as matérias marcadas para ver depois no banco de dados: ", err.Error())
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
			log.Error("Erro ao obter os dados das matérias marcadas para ver depois no banco de dados: ", err.Error())
			return nil, 0, err
		}
	}

	var articleList []article.Article
	for index, articleData := range articles {
		var articleDomain *article.Article
		if articleData.Proposition.Id != uuid.Nil {
			articleDomain, err = article.NewBuilder().
				Id(articleData.Id).
				Title(articleData.Proposition.Title).
				Content(articleData.Proposition.Content).
				AverageRating(articleData.AverageRating).
				NumberOfRatings(articleData.NumberOfRatings).
				UserRating(userArticles[index].Rating).
				ViewLater(userArticles[index].ViewLater).
				Type("Proposição").
				ReferenceDateTime(articleData.ReferenceDateTime).
				CreatedAt(articleData.CreatedAt).
				UpdatedAt(articleData.UpdatedAt).
				Build()
			if err != nil {
				log.Errorf("Erro ao validar os dados da matéria %s da proposição %s: %s", articleData.Id,
					articleData.Proposition.Id, err.Error())
				continue
			}
		} else {
			articleDomain, err = article.NewBuilder().
				Id(articleData.Id).
				Title(articleData.Newsletter.Title).
				Content(articleData.Newsletter.Description).
				AverageRating(articleData.AverageRating).
				NumberOfRatings(articleData.NumberOfRatings).
				UserRating(userArticles[index].Rating).
				ViewLater(userArticles[index].ViewLater).
				Type("Boletim").
				ReferenceDateTime(articleData.ReferenceDateTime).
				CreatedAt(articleData.CreatedAt).
				UpdatedAt(articleData.UpdatedAt).
				Build()
			if err != nil {
				log.Errorf("Erro ao validar os dados da matéria %s do boletim %s: %s", articleData.Id,
					articleData.Newsletter.Id, err.Error())
				continue
			}
		}

		articleList = append(articleList, *articleDomain)
	}

	var totalNumberOfArticles int
	if filter.Type == "Proposição" || filter.DeputyId != nil || filter.PartyId != nil || filter.ExternalAuthorId != nil {
		err = postgresConnection.Get(&totalNumberOfArticles, queries.UserArticle().Select().NumberOfPropositionsSavedToViewLater(),
			fmt.Sprint("%", filter.Content, "%"), &filter.DeputyId, &filter.PartyId, &filter.ExternalAuthorId,
			filter.StartDate, filter.EndDate, userId)
	} else if filter.Type == "Boletim" {
		err = postgresConnection.Get(&totalNumberOfArticles, queries.UserArticle().Select().NumberOfNewslettersSavedToViewLater(),
			fmt.Sprint("%", filter.Content, "%"), filter.StartDate, filter.EndDate, userId)
	} else {
		err = postgresConnection.Get(&totalNumberOfArticles, queries.UserArticle().Select().NumberOfArticlesSavedToViewLater(),
			fmt.Sprint("%", filter.Content, "%"), filter.StartDate, filter.EndDate, userId)
	}
	if err != nil {
		log.Error("Erro ao obter a quantidade total de matérias marcadas para ver depois no banco de dados: ", err.Error())
		return nil, 0, err
	}

	return articleList, totalNumberOfArticles, nil
}

func (instance Article) GetPropositionArticlesByNewsletterId(newsletterId uuid.UUID, userId uuid.UUID) ([]article.Article, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var propositionArticles []dto.Article
	err = postgresConnection.Select(&propositionArticles, queries.Article().Select().PropositionsByNewsletterId(),
		newsletterId)
	if err != nil {
		log.Errorf("Erro ao obter os dados das matérias das proposições do boletim %s no banco de dados: %s",
			newsletterId, err.Error())
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
			log.Error("Erro as buscar no banco de dados as matérias que o usuário avaliou e/ou salvou para ver depois: ", err.Error())
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

		articleDomain, err := articleBuilder.
			Id(articleData.Id).
			Title(articleData.Proposition.Title).
			Content(articleData.Proposition.Content).
			AverageRating(articleData.AverageRating).
			NumberOfRatings(articleData.NumberOfRatings).
			Type("Proposição").
			ReferenceDateTime(articleData.ReferenceDateTime).
			CreatedAt(articleData.CreatedAt).
			UpdatedAt(articleData.UpdatedAt).
			Build()
		if err != nil {
			log.Errorf("Erro ao validar os dados da matéria %s da proposição %s: %s", articleData.Id,
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
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var newsletterArticle dto.Article
	err = postgresConnection.Get(&newsletterArticle, queries.Article().Select().NewsletterByPropositionId(), propositionId)
	if err != nil {
		log.Errorf("Erro ao obter os dados do boletim pela matéria %s no banco de dados: %s", propositionId, err.Error())
		return nil, err
	}

	var userArticle dto.UserArticle
	if userId != uuid.Nil && newsletterArticle.Id != uuid.Nil {
		numberOfArticles := 1
		err = postgresConnection.Get(&userArticle, queries.UserArticle().Select().RatingsAndArticlesSavedForLaterViewing(
			numberOfArticles), userId, newsletterArticle.Id)
		if err != nil && err != sql.ErrNoRows {
			log.Errorf("Erro as buscar no banco de dados as informações da matéria %s que o usuário %s avaliou e/ou "+
				"salvou para ver depois: %s", newsletterArticle.Id, userId, err.Error())
			return nil, err
		}
	}

	articleBuilder := article.NewBuilder()

	if userArticle.Article.Id != uuid.Nil {
		articleBuilder.UserRating(userArticle.Rating).ViewLater(userArticle.ViewLater)
	}

	articleDomain, err := articleBuilder.
		Id(newsletterArticle.Id).
		Title(newsletterArticle.Newsletter.Title).
		Content(newsletterArticle.Newsletter.Description).
		AverageRating(newsletterArticle.AverageRating).
		NumberOfRatings(newsletterArticle.NumberOfRatings).
		Type("Boletim").
		ReferenceDateTime(newsletterArticle.ReferenceDateTime).
		CreatedAt(newsletterArticle.CreatedAt).
		UpdatedAt(newsletterArticle.UpdatedAt).
		Build()
	if err != nil {
		log.Errorf("Erro ao validar os dados da matéria %s do boletim %s: %s", newsletterArticle.Id,
			newsletterArticle.Newsletter.Id, err.Error())
		return nil, err
	}

	return articleDomain, nil
}

func (instance Article) SaveArticleRating(userId uuid.UUID, articleId uuid.UUID, rating int) error {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		return err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	sqlResult, err := postgresConnection.Exec(queries.UserArticle().Update().Rating(), rating, userId, articleId)
	if err != nil {
		log.Errorf("Erro ao atualizar avaliação do artigo %s com o usuário %s: %s", articleId, userId, err)
		return err
	}

	rowsAffected, err := sqlResult.RowsAffected()
	if err == nil && rowsAffected == 0 {
		_, err = postgresConnection.Exec(queries.UserArticle().Insert().Rating(), userId, articleId, rating)
		if err != nil {
			log.Errorf("Erro ao inserir avaliação do artigo %s com o usuário %s: %s", articleId, userId, err)
			return err
		}
	} else if err != nil {
		log.Errorf("Erro ao extrair a quantidade de linhas afetadas pela atualização da avaliação do artigo %s "+
			"com o usuário %s: %s", articleId, userId, err)
		return err
	}

	return nil
}

func (instance Article) SaveArticleToViewLater(userId uuid.UUID, articleId uuid.UUID, viewLater bool) error {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		return err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	sqlResult, err := postgresConnection.Exec(queries.UserArticle().Update().ViewLater(), viewLater, userId, articleId)
	if err != nil {
		log.Errorf("Erro ao atualizar marcação de ver depois do artigo %s com o usuário %s: %s",
			articleId, userId, err)
		return err
	}

	rowsAffected, err := sqlResult.RowsAffected()
	if err == nil && rowsAffected == 0 {
		_, err = postgresConnection.Exec(queries.UserArticle().Insert().ViewLater(), userId, articleId, viewLater)
		if err != nil {
			log.Errorf("Erro ao inserir marcação de ver depois do artigo %s com o usuário %s: %s",
				articleId, userId, err)
			return err
		}
	} else if err != nil {
		log.Errorf("Erro ao extrair a quantidade de linhas afetadas pela atualização da marcação de ver depois "+
			"do artigo %s com o usuário %s: %s", articleId, userId, err)
		return err
	}

	return nil
}
