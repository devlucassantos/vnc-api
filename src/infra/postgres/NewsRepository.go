package postgres

import (
	"github.com/labstack/gommon/log"
	"vnc-read-api/api/endpoints/dto/filter"
	"vnc-read-api/core/domains/news"
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

func (instance News) GetNews(filter filter.NewsFilter) ([]news.News, int, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		return nil, 0, err
	}
	defer instance.connectionManager.endConnection(postgresConnection)

	var newsList []dto.News
	err = postgresConnection.Select(&newsList, queries.News().Select().All(),
		filter.PaginationFilter.CalculateOffset(), filter.PaginationFilter.GetItemsPerPage())
	if err != nil {
		log.Error("Erro ao obter os dados das matérias no banco de dados: ", err.Error())
		return nil, 0, err
	}

	var newsDomainList []news.News
	for _, newsData := range newsList {
		var newsDomain *news.News
		if newsData.Proposition.Id.ID() != 0 {
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
	err = postgresConnection.Get(&totalNumberOfNews, queries.News().Select().TotalNumber())
	if err != nil {
		log.Error("Erro ao obter a quantidade total de matérias no banco de dados: ", err.Error())
		return nil, 0, err
	}

	return newsDomainList, totalNumberOfNews, nil
}
