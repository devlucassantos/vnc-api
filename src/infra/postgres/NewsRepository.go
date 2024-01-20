package postgres

import (
	"fmt"
	"github.com/devlucassantos/vnc-domains/src/domains/news"
	"github.com/labstack/gommon/log"
	"vnc-read-api/core/filters"
	"vnc-read-api/infra/dto"
	"vnc-read-api/infra/postgres/queries"
)

type News struct {
	connectionManager ConnectionManagerInterface
}

func NewNewsRepository(connectionManager ConnectionManagerInterface) *News {
	return &News{
		connectionManager: connectionManager,
	}
}

func (instance News) GetNews(filter filters.NewsFilter) ([]news.News, int, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		return nil, 0, err
	}
	defer instance.connectionManager.endConnection(postgresConnection)

	var newsList []dto.News
	if filter.Type == "Proposição" || filter.DeputyId != nil || filter.PartyId != nil || filter.OrganizationId != nil {
		err = postgresConnection.Select(&newsList, queries.News().Select().Propositions(),
			fmt.Sprint("%", filter.Content, "%"), filter.DeputyId, filter.PartyId, filter.OrganizationId,
			filter.StartDate, filter.EndDate, filter.PaginationFilter.CalculateOffset(), filter.PaginationFilter.GetItemsPerPage())
	} else if filter.Type == "Boletim" {
		err = postgresConnection.Select(&newsList, queries.News().Select().Newsletters(),
			fmt.Sprint("%", filter.Content, "%"), filter.StartDate, filter.EndDate,
			filter.PaginationFilter.CalculateOffset(), filter.PaginationFilter.GetItemsPerPage())
	} else {
		err = postgresConnection.Select(&newsList, queries.News().Select().All(),
			fmt.Sprint("%", filter.Content, "%"), filter.StartDate, filter.EndDate,
			filter.PaginationFilter.CalculateOffset(), filter.PaginationFilter.GetItemsPerPage())
	}
	if err != nil {
		log.Error("Erro ao obter os dados das matérias no banco de dados: ", err.Error())
		return nil, 0, err
	}

	var newsDomainList []news.News
	for _, newsData := range newsList {
		var newsDomain *news.News
		if newsData.Proposition != nil && newsData.Proposition.Id.ID() != 0 {
			newsDomain, err = news.NewBuilder().
				Id(newsData.Proposition.Id).
				Title(newsData.Proposition.Title).
				Content(newsData.Proposition.Content).
				Type("Proposição").
				Active(newsData.Proposition.Active).
				CreatedAt(newsData.Proposition.CreatedAt).
				UpdatedAt(newsData.Proposition.UpdatedAt).
				Build()
			if err != nil {
				log.Errorf("Erro construindo a estrutura de dados da matéria da proposição %s: %s",
					newsData.Proposition.Id, err.Error())
				continue
			}
		} else {
			newsDomain, err = news.NewBuilder().
				Id(newsData.Newsletter.Id).
				Title(newsData.Newsletter.Title).
				Content(newsData.Newsletter.Content).
				Type("Boletim").
				Active(newsData.Newsletter.Active).
				CreatedAt(newsData.Newsletter.CreatedAt).
				UpdatedAt(newsData.Newsletter.UpdatedAt).
				Build()
			if err != nil {
				log.Errorf("Erro construindo a estrutura de dados da matéria do boletim %s: %s",
					newsData.Newsletter.Id, err.Error())
				continue
			}
		}

		newsDomainList = append(newsDomainList, *newsDomain)
	}

	var totalNumberOfNews int
	if filter.Type == "Proposição" || filter.DeputyId != nil || filter.PartyId != nil || filter.OrganizationId != nil {
		err = postgresConnection.Get(&totalNumberOfNews, queries.News().Select().TotalNumberOfPropositions(),
			fmt.Sprint("%", filter.Content, "%"), &filter.DeputyId, &filter.PartyId, &filter.OrganizationId,
			filter.StartDate, filter.EndDate)
	} else if filter.Type == "Boletim" {
		err = postgresConnection.Get(&totalNumberOfNews, queries.News().Select().TotalNumberOfNewsletters(),
			fmt.Sprint("%", filter.Content, "%"), filter.StartDate, filter.EndDate)
	} else {
		err = postgresConnection.Get(&totalNumberOfNews, queries.News().Select().TotalNumberOfNews(),
			fmt.Sprint("%", filter.Content, "%"), filter.StartDate, filter.EndDate)
	}
	if err != nil {
		log.Error("Erro ao obter a quantidade total de matérias no banco de dados: ", err.Error())
		return nil, 0, err
	}

	return newsDomainList, totalNumberOfNews, nil
}

func (instance News) GetTrending(filter filters.NewsFilter) ([]news.News, int, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		return nil, 0, err
	}
	defer instance.connectionManager.endConnection(postgresConnection)

	var newsList []dto.News
	if filter.Type == "Proposição" || filter.DeputyId != nil || filter.PartyId != nil || filter.OrganizationId != nil {
		err = postgresConnection.Select(&newsList, queries.News().Select().TrendingPropositions(), fmt.Sprint("%",
			filter.Content, "%"), filter.DeputyId, filter.PartyId, filter.OrganizationId, filter.StartDate, filter.EndDate,
			filter.PaginationFilter.CalculateOffset(), filter.PaginationFilter.GetItemsPerPage())
	} else if filter.Type == "Boletim" {
		err = postgresConnection.Select(&newsList, queries.News().Select().TrendingNewsletters(),
			fmt.Sprint("%", filter.Content, "%"), filter.StartDate, filter.EndDate,
			filter.PaginationFilter.CalculateOffset(), filter.PaginationFilter.GetItemsPerPage())
	} else {
		err = postgresConnection.Select(&newsList, queries.News().Select().TrendingNews(),
			fmt.Sprint("%", filter.Content, "%"), filter.StartDate, filter.EndDate,
			filter.PaginationFilter.CalculateOffset(), filter.PaginationFilter.GetItemsPerPage())
	}
	if err != nil {
		log.Error("Erro ao obter os dados das matérias no banco de dados: ", err.Error())
		return nil, 0, err
	}

	var newsDomainList []news.News
	for _, newsData := range newsList {
		var newsDomain *news.News
		if newsData.Proposition != nil && newsData.Proposition.Id.ID() != 0 {
			newsDomain, err = news.NewBuilder().
				Id(newsData.Proposition.Id).
				Title(newsData.Proposition.Title).
				Content(newsData.Proposition.Content).
				Type("Proposição").
				Active(newsData.Proposition.Active).
				CreatedAt(newsData.Proposition.CreatedAt).
				UpdatedAt(newsData.Proposition.UpdatedAt).
				Build()
			if err != nil {
				log.Errorf("Erro construindo a estrutura de dados da matéria da proposição %s: %s",
					newsData.Proposition.Id, err.Error())
				continue
			}
		} else {
			newsDomain, err = news.NewBuilder().
				Id(newsData.Newsletter.Id).
				Title(newsData.Newsletter.Title).
				Content(newsData.Newsletter.Content).
				Type("Boletim").
				Active(newsData.Newsletter.Active).
				CreatedAt(newsData.Newsletter.CreatedAt).
				UpdatedAt(newsData.Newsletter.UpdatedAt).
				Build()
			if err != nil {
				log.Errorf("Erro construindo a estrutura de dados da matéria do boletim %s: %s",
					newsData.Newsletter.Id, err.Error())
				continue
			}
		}

		newsDomainList = append(newsDomainList, *newsDomain)
	}

	var totalNumberOfNews int
	if filter.Type == "Proposição" || filter.DeputyId != nil || filter.PartyId != nil || filter.OrganizationId != nil {
		err = postgresConnection.Get(&totalNumberOfNews, queries.News().Select().TotalNumberOfPropositions(),
			fmt.Sprint("%", filter.Content, "%"), &filter.DeputyId, &filter.PartyId, &filter.OrganizationId,
			filter.StartDate, filter.EndDate)
	} else if filter.Type == "Boletim" {
		err = postgresConnection.Get(&totalNumberOfNews, queries.News().Select().TotalNumberOfNewsletters(),
			fmt.Sprint("%", filter.Content, "%"), filter.StartDate, filter.EndDate)
	} else {
		err = postgresConnection.Get(&totalNumberOfNews, queries.News().Select().TotalNumberOfNews(),
			fmt.Sprint("%", filter.Content, "%"), filter.StartDate, filter.EndDate)
	}
	if err != nil {
		log.Error("Erro ao obter a quantidade total de matérias no banco de dados: ", err.Error())
		return nil, 0, err
	}

	return newsDomainList, totalNumberOfNews, nil
}
