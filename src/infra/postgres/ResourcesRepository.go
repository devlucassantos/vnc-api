package postgres

import (
	"fmt"
	"github.com/devlucassantos/vnc-domains/src/domains/articletype"
	"github.com/devlucassantos/vnc-domains/src/domains/deputy"
	"github.com/devlucassantos/vnc-domains/src/domains/eventsituation"
	"github.com/devlucassantos/vnc-domains/src/domains/eventtype"
	"github.com/devlucassantos/vnc-domains/src/domains/externalauthor"
	"github.com/devlucassantos/vnc-domains/src/domains/externalauthortype"
	"github.com/devlucassantos/vnc-domains/src/domains/legislativebody"
	"github.com/devlucassantos/vnc-domains/src/domains/legislativebodytype"
	"github.com/devlucassantos/vnc-domains/src/domains/party"
	"github.com/devlucassantos/vnc-domains/src/domains/propositiontype"
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

func (instance Resources) GetArticleTypes() ([]articletype.ArticleType, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("Error creating a connection to the Postgres database: ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var articleTypesData []dto.ArticleType
	err = postgresConnection.Select(&articleTypesData, queries.ArticleType().Select().All())
	if err != nil {
		log.Error("Error retrieving article types data from the database: ", err.Error())
		return nil, err
	}

	var articleTypes []articletype.ArticleType
	for _, articleType := range articleTypesData {
		articleTypeDomain, err := articletype.NewBuilder().
			Id(articleType.Id).
			Description(articleType.Description).
			Codes(articleType.Codes).
			Color(articleType.Color).
			Build()
		if err != nil {
			log.Errorf("Error validating data for article type %s: %s", articleType.Id, err.Error())
			return nil, err
		}
		articleTypes = append(articleTypes, *articleTypeDomain)
	}

	return articleTypes, nil
}

func (instance Resources) GetPropositionTypes() ([]propositiontype.PropositionType, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("Error creating a connection to the Postgres database: ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var propositionTypesData []dto.PropositionType
	err = postgresConnection.Select(&propositionTypesData, queries.PropositionType().Select().All())
	if err != nil {
		log.Error("Error retrieving proposition types data from the database: ", err.Error())
		return nil, err
	}

	var propositionTypes []propositiontype.PropositionType
	for _, propositionType := range propositionTypesData {
		propositionTypeDomain, err := propositiontype.NewBuilder().
			Id(propositionType.Id).
			Description(propositionType.Description).
			Color(propositionType.Color).
			Build()
		if err != nil {
			log.Errorf("Error validating data for proposition type %s: %s", propositionType.Id, err.Error())
			return nil, err
		}
		propositionTypes = append(propositionTypes, *propositionTypeDomain)
	}

	return propositionTypes, nil
}

func (instance Resources) GetParties() ([]party.Party, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("Error creating a connection to the Postgres database: ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var partyDtos []dto.Party
	err = postgresConnection.Select(&partyDtos, queries.Party().Select().All())
	if err != nil {
		log.Error("Error retrieving party data from the database: ", err.Error())
		return nil, err
	}

	var parties []party.Party
	for _, partyData := range partyDtos {
		partyDomain, err := party.NewBuilder().
			Id(partyData.Id).
			Name(partyData.Name).
			Acronym(partyData.Acronym).
			ImageUrl(partyData.ImageUrl).
			ImageDescription(fmt.Sprintf("Logo do %s (%s)", partyData.Name, partyData.Acronym)).
			Build()
		if err != nil {
			log.Errorf("Error validating data for party %s: %s", partyData.Id, err.Error())
			return nil, err
		}
		parties = append(parties, *partyDomain)
	}

	return parties, nil
}

func (instance Resources) GetDeputies() ([]deputy.Deputy, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("Error creating a connection to the Postgres database: ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var deputiesData []dto.Deputy
	err = postgresConnection.Select(&deputiesData, queries.Deputy().Select().All())
	if err != nil {
		log.Error("Error retrieving deputies data from the database: ", err.Error())
		return nil, err
	}

	var deputies []deputy.Deputy
	for _, deputyData := range deputiesData {
		currentParty, err := party.NewBuilder().
			Id(deputyData.Party.Id).
			Name(deputyData.Party.Name).
			Acronym(deputyData.Party.Acronym).
			ImageUrl(deputyData.Party.ImageUrl).
			ImageDescription(fmt.Sprintf("Logo do %s (%s)", deputyData.Party.Name, deputyData.Party.Acronym)).
			Build()
		if err != nil {
			log.Errorf("Error validating data for the current party %s of deputy %s: %s", deputyData.Party.Id,
				deputyData.Id, err.Error())
			return nil, err
		}

		deputyDomain, err := deputy.NewBuilder().
			Id(deputyData.Id).
			Name(deputyData.Name).
			ElectoralName(deputyData.ElectoralName).
			ImageUrl(deputyData.ImageUrl).
			ImageDescription(fmt.Sprintf("Foto do(a) deputado(a) federal %s (%s-%s)", deputyData.Name,
				deputyData.Party.Acronym, deputyData.FederatedUnit)).
			Party(*currentParty).
			FederatedUnit(deputyData.FederatedUnit).
			Build()
		if err != nil {
			log.Errorf("Error validating data for deputy %s: %s", deputyData.Id, err.Error())
			return nil, err
		}
		deputies = append(deputies, *deputyDomain)
	}

	return deputies, nil
}

func (instance Resources) GetExternalAuthors() ([]externalauthor.ExternalAuthor, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("Error creating a connection to the Postgres database: ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var externalAuthorsData []dto.ExternalAuthor
	err = postgresConnection.Select(&externalAuthorsData, queries.ExternalAuthor().Select().All())
	if err != nil {
		log.Error("Error retrieving external authors data from the database: ", err.Error())
		return nil, err
	}

	var externalAuthors []externalauthor.ExternalAuthor
	for _, externalAuthor := range externalAuthorsData {
		externalAuthorType, err := externalauthortype.NewBuilder().
			Id(externalAuthor.ExternalAuthorType.Id).
			Description(externalAuthor.ExternalAuthorType.Description).
			Build()
		if err != nil {
			log.Errorf("Error validating data for external author type %s for external author %s: %s",
				externalAuthor.ExternalAuthorType.Id, externalAuthor.Id, err.Error())
			return nil, err
		}

		externalAuthorDomain, err := externalauthor.NewBuilder().
			Id(externalAuthor.Id).
			Name(externalAuthor.Name).
			Type(*externalAuthorType).
			Build()
		if err != nil {
			log.Errorf("Error validating data for external author %s: %s", externalAuthor.Id, err.Error())
			return nil, err
		}
		externalAuthors = append(externalAuthors, *externalAuthorDomain)
	}

	return externalAuthors, nil
}

func (instance Resources) GetLegislativeBodies() ([]legislativebody.LegislativeBody, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("Error creating a connection to the Postgres database: ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var legislativeBodiesData []dto.LegislativeBody
	err = postgresConnection.Select(&legislativeBodiesData, queries.LegislativeBody().Select().All())
	if err != nil {
		log.Error("Error retrieving legislative bodies data from the database: ", err.Error())
		return nil, err
	}

	var legislativeBodies []legislativebody.LegislativeBody
	for _, legislativeBody := range legislativeBodiesData {
		legislativeBodyType, err := legislativebodytype.NewBuilder().
			Id(legislativeBody.LegislativeBodyType.Id).
			Description(legislativeBody.LegislativeBodyType.Description).
			Build()
		if err != nil {
			log.Errorf("Error validating data for the legislative body type %s of legislative body %s: %s",
				legislativeBody.LegislativeBodyType.Id, legislativeBody.Id, err.Error())
			return nil, err
		}

		legislativeBodyDomain, err := legislativebody.NewBuilder().
			Id(legislativeBody.Id).
			Name(legislativeBody.Name).
			Acronym(legislativeBody.Acronym).
			Type(*legislativeBodyType).
			Build()
		if err != nil {
			log.Errorf("Error validating data for legislative body %s: %s", legislativeBody.Id, err.Error())
			return nil, err
		}
		legislativeBodies = append(legislativeBodies, *legislativeBodyDomain)
	}

	return legislativeBodies, nil
}

func (instance Resources) GetEventTypes() ([]eventtype.EventType, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("Error creating a connection to the Postgres database: ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var eventTypesData []dto.EventType
	err = postgresConnection.Select(&eventTypesData, queries.EventType().Select().All())
	if err != nil {
		log.Error("Error retrieving event types data from the database: ", err.Error())
		return nil, err
	}

	var eventTypes []eventtype.EventType
	for _, eventType := range eventTypesData {
		eventTypeDomain, err := eventtype.NewBuilder().
			Id(eventType.Id).
			Description(eventType.Description).
			Color(eventType.Color).
			Build()
		if err != nil {
			log.Errorf("Error validating data for event type %s: %s", eventType.Id, err.Error())
			return nil, err
		}
		eventTypes = append(eventTypes, *eventTypeDomain)
	}

	return eventTypes, nil
}

func (instance Resources) GetEventSituations() ([]eventsituation.EventSituation, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("Error creating a connection to the Postgres database: ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var eventSituationsData []dto.EventSituation
	err = postgresConnection.Select(&eventSituationsData, queries.EventSituation().Select().All())
	if err != nil {
		log.Error("Error retrieving event situations data from the database: ", err.Error())
		return nil, err
	}

	var eventSituations []eventsituation.EventSituation
	for _, eventSituation := range eventSituationsData {
		eventSituationDomain, err := eventsituation.NewBuilder().
			Id(eventSituation.Id).
			Description(eventSituation.Description).
			Color(eventSituation.Color).
			Build()
		if err != nil {
			log.Errorf("Error validating data for event situation %s: %s", eventSituation.Id, err.Error())
			return nil, err
		}
		eventSituations = append(eventSituations, *eventSituationDomain)
	}

	return eventSituations, nil
}
