package services

import (
	"github.com/labstack/gommon/log"
	"vnc-read-api/core/domains/deputy"
	"vnc-read-api/core/domains/organization"
	"vnc-read-api/core/domains/party"
	"vnc-read-api/core/interfaces/repositories"
)

type Resources struct {
	repository repositories.Resources
}

func NewResourcesService(repository repositories.Resources) *Resources {
	return &Resources{
		repository: repository,
	}
}

func (instance Resources) GetResources() ([]party.Party, []deputy.Deputy, []organization.Organization, error) {
	parties, err := instance.repository.GetParties()
	if err != nil {
		log.Errorf("Erro ao obter os dados dos partidos no banco de dados: %s", err.Error())
		return nil, nil, nil, err
	}

	deputies, err := instance.repository.GetDeputies()
	if err != nil {
		log.Errorf("Erro ao obter os dados dos deputados no banco de dados: %s", err.Error())
		return nil, nil, nil, err
	}

	organizations, err := instance.repository.GetOrganizations()
	if err != nil {
		log.Errorf("Erro ao obter os dados das organizações no banco de dados: %s", err.Error())
		return nil, nil, nil, err
	}

	return parties, deputies, organizations, nil
}
