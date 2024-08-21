package services

import (
	"github.com/devlucassantos/vnc-domains/src/domains/deputy"
	"github.com/devlucassantos/vnc-domains/src/domains/external"
	"github.com/devlucassantos/vnc-domains/src/domains/party"
	"github.com/devlucassantos/vnc-domains/src/domains/proptype"
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

func (instance Resources) GetResources() ([]proptype.PropositionType, []party.Party, []deputy.Deputy,
	[]external.ExternalAuthor, error) {
	propositionTypes, err := instance.repository.GetPropositionTypes(nil)
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

	return propositionTypes, parties, deputies, externalAuthors, nil
}

func (instance Resources) GetPropositionTypes(propositionTypeIds []uuid.UUID) ([]proptype.PropositionType, error) {
	return instance.repository.GetPropositionTypes(propositionTypeIds)
}
