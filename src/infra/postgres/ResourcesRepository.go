package postgres

import (
	"github.com/devlucassantos/vnc-domains/src/domains/deputy"
	"github.com/devlucassantos/vnc-domains/src/domains/organization"
	"github.com/devlucassantos/vnc-domains/src/domains/party"
	"github.com/labstack/gommon/log"
	"vnc-read-api/infra/dto"
	"vnc-read-api/infra/postgres/queries"
)

type Resources struct {
	connectionManager ConnectionManagerInterface
}

func NewResourcesRepository(connectionManager ConnectionManagerInterface) *Resources {
	return &Resources{
		connectionManager: connectionManager,
	}
}

func (instance Resources) GetParties() ([]party.Party, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		return nil, err
	}
	defer instance.connectionManager.endConnection(postgresConnection)

	var partiesData []dto.Party
	err = postgresConnection.Select(&partiesData, queries.Party().Select().All())
	if err != nil {
		log.Errorf("Erro ao obter os dados dos partidos no banco de dados: %s", err.Error())
		return nil, err
	}

	var parties []party.Party
	for _, partyData := range partiesData {
		partyDomain, err := party.NewBuilder().
			Id(partyData.Id).
			Code(partyData.Code).
			Name(partyData.Name).
			Acronym(partyData.Acronym).
			ImageUrl(partyData.ImageUrl).
			Active(partyData.Active).
			CreatedAt(partyData.CreatedAt).
			UpdatedAt(partyData.UpdatedAt).
			Build()
		if err != nil {
			log.Errorf("Erro construindo a estrutura de dados do partido %s: %s", partyDomain.Id, err.Error())
		}

		parties = append(parties, *partyDomain)
	}

	return parties, nil
}

func (instance Resources) GetDeputies() ([]deputy.Deputy, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		return nil, err
	}
	defer instance.connectionManager.endConnection(postgresConnection)

	var deputiesData []dto.Deputy
	err = postgresConnection.Select(&deputiesData, queries.Deputy().Select().All())
	if err != nil {
		log.Errorf("Erro ao obter os dados dos deputados no banco de dados: %s", err.Error())
		return nil, err
	}

	var deputies []deputy.Deputy
	for _, deputyData := range deputiesData {
		currentParty, err := party.NewBuilder().
			Id(deputyData.Party.Id).
			Code(deputyData.Party.Code).
			Name(deputyData.Party.Name).
			Acronym(deputyData.Party.Acronym).
			ImageUrl(deputyData.Party.ImageUrl).
			Active(deputyData.Party.Active).
			CreatedAt(deputyData.Party.CreatedAt).
			UpdatedAt(deputyData.Party.UpdatedAt).
			Build()
		if err != nil {
			log.Errorf("Erro construindo a estrutura de dados do partido atual %s do(a) deputado(a) %s: %s",
				deputyData.Party.Id, deputyData.Id, err.Error())
		}

		deputyDomain, err := deputy.NewBuilder().
			Id(deputyData.Id).
			Code(deputyData.Code).
			Cpf(deputyData.Cpf).
			Name(deputyData.Name).
			ElectoralName(deputyData.ElectoralName).
			ImageUrl(deputyData.ImageUrl).
			Party(*currentParty).
			Active(deputyData.Active).
			CreatedAt(deputyData.CreatedAt).
			UpdatedAt(deputyData.UpdatedAt).
			Build()
		if err != nil {
			log.Errorf("Erro construindo a estrutura de dados do(a) deputado(a) %s: %s", deputyData.Id, err.Error())
			continue
		}

		deputies = append(deputies, *deputyDomain)
	}

	return deputies, nil
}

func (instance Resources) GetOrganizations() ([]organization.Organization, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		return nil, err
	}
	defer instance.connectionManager.endConnection(postgresConnection)

	var organizationsData []dto.Organization
	err = postgresConnection.Select(&organizationsData, queries.Organization().Select().All())
	if err != nil {
		log.Errorf("Erro ao obter os dados das organizações no banco de dados: %s", err.Error())
		return nil, err
	}

	var organizations []organization.Organization
	for _, organizationData := range organizationsData {
		organizationBuilder := organization.NewBuilder().
			Id(organizationData.Id)

		if organizationData.Code > 0 {
			organizationBuilder.Code(organizationData.Code)
		}

		organizationDomain, err := organizationBuilder.
			Name(organizationData.Name).
			Acronym(organizationData.Acronym).
			Nickname(organizationData.Nickname).
			Type(organizationData.Type).
			Active(organizationData.Active).
			CreatedAt(organizationData.CreatedAt).
			UpdatedAt(organizationData.UpdatedAt).
			Build()
		if err != nil {
			log.Errorf("Erro construindo a estrutura de dados do partido %s: %s", organizationDomain.Id, err.Error())
		}

		organizations = append(organizations, *organizationDomain)
	}

	return organizations, nil
}
