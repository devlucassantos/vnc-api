package services

import (
	"github.com/devlucassantos/vnc-domains/src/domains/articletype"
	"github.com/devlucassantos/vnc-domains/src/domains/deputy"
	"github.com/devlucassantos/vnc-domains/src/domains/external"
	"github.com/devlucassantos/vnc-domains/src/domains/party"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"vnc-api/core/interfaces/repositories"
)

type Resources struct {
	repository repositories.Resources
}

func NewResourcesService(repository repositories.Resources) *Resources {
	return &Resources{
		repository: repository,
	}
}

func (instance Resources) GetResources() ([]articletype.ArticleType, []party.Party, []deputy.Deputy,
	[]external.ExternalAuthor, error) {
	articleTypes, err := instance.repository.GetArticleTypes(nil)
	if err != nil {
		log.Errorf("Erro ao obter os dados dos tipos de proposição no banco de dados: %s", err.Error())
		return nil, nil, nil, nil, err
	}

	parties, err := instance.repository.GetParties()
	if err != nil {
		log.Errorf("Erro ao obter os dados dos partidos no banco de dados: %s", err.Error())
		return nil, nil, nil, nil, err
	}

	deputies, err := instance.repository.GetDeputies()
	if err != nil {
		log.Errorf("Erro ao obter os dados dos deputados no banco de dados: %s", err.Error())
		return nil, nil, nil, nil, err
	}

	externalAuthors, err := instance.repository.GetExternalAuthors()
	if err != nil {
		log.Errorf("Erro ao obter os dados dos autores externos no banco de dados: %s", err.Error())
		return nil, nil, nil, nil, err
	}

	return articleTypes, parties, deputies, externalAuthors, nil
}

func (instance Resources) GetArticleTypes(articleTypeIds []uuid.UUID) ([]articletype.ArticleType, error) {
	return instance.repository.GetArticleTypes(articleTypeIds)
}
