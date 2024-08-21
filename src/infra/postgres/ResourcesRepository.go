package postgres

import (
	"github.com/devlucassantos/vnc-domains/src/domains/deputy"
	"github.com/devlucassantos/vnc-domains/src/domains/external"
	"github.com/devlucassantos/vnc-domains/src/domains/party"
	"github.com/devlucassantos/vnc-domains/src/domains/proptype"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"vnc-api/infra/dto"
	"vnc-api/infra/postgres/queries"
)

type Resources struct {
	connectionManager connectionManagerInterface
}

func NewResourcesRepository(connectionManager connectionManagerInterface) *Resources {
	return &Resources{
		connectionManager: connectionManager,
	}
}

func (instance Resources) GetPropositionTypes(propositionTypeIds []uuid.UUID) ([]proptype.PropositionType, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("Erro ao tentar se conectar com o Postgres: ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var propositionTypeIdsAsInterfaceSlice []interface{}
	for _, propositionTypeId := range propositionTypeIds {
		propositionTypeIdsAsInterfaceSlice = append(propositionTypeIdsAsInterfaceSlice, propositionTypeId)
	}

	var propositionTypeData []dto.PropositionType
	if propositionTypeIdsAsInterfaceSlice != nil {
		err = postgresConnection.Select(&propositionTypeData, queries.PropositionType().Select().In(
			len(propositionTypeIdsAsInterfaceSlice)), propositionTypeIdsAsInterfaceSlice...)
	} else {
		err = postgresConnection.Select(&propositionTypeData, queries.PropositionType().Select().All())
	}
	if err != nil {
		log.Errorf("Erro ao obter os dados dos tipos de proposição no banco de dados: %s", err.Error())
		return nil, err
	}

	var propositionTypes []proptype.PropositionType
	for _, propositionType := range propositionTypeData {
		propositionTypeDomain, err := proptype.NewBuilder().
			Id(propositionType.Id).
			Description(propositionType.Description).
			Color(propositionType.Color).
			SortOrder(propositionType.SortOrder).
			CreatedAt(propositionType.CreatedAt).
			UpdatedAt(propositionType.UpdatedAt).
			Build()
		if err != nil {
			log.Errorf("Erro ao validar os dados do tipo de proposição %s: %s", propositionType.Id, err.Error())
		}

		propositionTypes = append(propositionTypes, *propositionTypeDomain)
	}

	return propositionTypes, nil
}

func (instance Resources) GetParties() ([]party.Party, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("Erro ao tentar se conectar com o Postgres: ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

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
			Name(partyData.Name).
			Acronym(partyData.Acronym).
			ImageUrl(partyData.ImageUrl).
			CreatedAt(partyData.CreatedAt).
			UpdatedAt(partyData.UpdatedAt).
			Build()
		if err != nil {
			log.Errorf("Erro ao validar os dados do partido %s: %s", partyDomain.Id, err.Error())
		}

		parties = append(parties, *partyDomain)
	}

	return parties, nil
}

func (instance Resources) GetDeputies() ([]deputy.Deputy, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("Erro ao tentar se conectar com o Postgres: ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

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
			Name(deputyData.Party.Name).
			Acronym(deputyData.Party.Acronym).
			ImageUrl(deputyData.Party.ImageUrl).
			CreatedAt(deputyData.Party.CreatedAt).
			UpdatedAt(deputyData.Party.UpdatedAt).
			Build()
		if err != nil {
			log.Errorf("Erro ao validar os dados do partido atual %s do(a) deputado(a) %s: %s",
				deputyData.Party.Id, deputyData.Id, err.Error())
		}

		deputyDomain, err := deputy.NewBuilder().
			Id(deputyData.Id).
			Name(deputyData.Name).
			ElectoralName(deputyData.ElectoralName).
			ImageUrl(deputyData.ImageUrl).
			Party(*currentParty).
			CreatedAt(deputyData.CreatedAt).
			UpdatedAt(deputyData.UpdatedAt).
			Build()
		if err != nil {
			log.Errorf("Erro ao validar os dados do(a) deputado(a) %s: %s", deputyData.Id, err.Error())
			continue
		}

		deputies = append(deputies, *deputyDomain)
	}

	return deputies, nil
}

func (instance Resources) GetExternalAuthors() ([]external.ExternalAuthor, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("Erro ao tentar se conectar com o Postgres: ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var externalAuthorsData []dto.ExternalAuthor
	err = postgresConnection.Select(&externalAuthorsData, queries.ExternalAuthor().Select().All())
	if err != nil {
		log.Errorf("Erro ao obter os dados dos autores externos no banco de dados: %s", err.Error())
		return nil, err
	}

	var externalAuthors []external.ExternalAuthor
	for _, externalAuthorData := range externalAuthorsData {
		externalAuthorDomain, err := external.NewBuilder().
			Id(externalAuthorData.Id).
			Name(externalAuthorData.Name).
			Type(externalAuthorData.Type).
			CreatedAt(externalAuthorData.CreatedAt).
			UpdatedAt(externalAuthorData.UpdatedAt).
			Build()
		if err != nil {
			log.Errorf("Erro ao validar os dados do autor externo %s: %s", externalAuthorDomain.Id,
				err.Error())
		}

		externalAuthors = append(externalAuthors, *externalAuthorDomain)
	}

	return externalAuthors, nil
}
