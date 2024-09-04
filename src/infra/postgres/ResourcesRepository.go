package postgres

import (
	"github.com/devlucassantos/vnc-domains/src/domains/articletype"
	"github.com/devlucassantos/vnc-domains/src/domains/deputy"
	"github.com/devlucassantos/vnc-domains/src/domains/external"
	"github.com/devlucassantos/vnc-domains/src/domains/party"
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

func (instance Resources) GetArticleTypes(articleTypeIds []uuid.UUID) ([]articletype.ArticleType, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("Erro ao tentar se conectar com o Postgres: ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var articleTypeIdsAsInterfaceSlice []interface{}
	for _, articleTypeId := range articleTypeIds {
		articleTypeIdsAsInterfaceSlice = append(articleTypeIdsAsInterfaceSlice, articleTypeId)
	}

	var articleTypeData []dto.ArticleType
	if articleTypeIdsAsInterfaceSlice != nil {
		err = postgresConnection.Select(&articleTypeData, queries.ArticleType().Select().In(
			len(articleTypeIdsAsInterfaceSlice)), articleTypeIdsAsInterfaceSlice...)
	} else {
		err = postgresConnection.Select(&articleTypeData, queries.ArticleType().Select().All())
	}
	if err != nil {
		log.Errorf("Erro ao obter os dados dos tipos de matéria no banco de dados: %s", err.Error())
		return nil, err
	}

	var articleTypes []articletype.ArticleType
	for _, articleType := range articleTypeData {
		articleTypeDomain, err := articletype.NewBuilder().
			Id(articleType.Id).
			Description(articleType.Description).
			Color(articleType.Color).
			SortOrder(articleType.SortOrder).
			CreatedAt(articleType.CreatedAt).
			UpdatedAt(articleType.UpdatedAt).
			Build()
		if err != nil {
			log.Errorf("Erro ao validar os dados do tipo de matéria %s: %s", articleType.Id, err.Error())
		}

		articleTypes = append(articleTypes, *articleTypeDomain)
	}

	return articleTypes, nil
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
