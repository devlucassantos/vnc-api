package services

import (
	"github.com/devlucassantos/vnc-domains/src/domains/articletype"
	"github.com/devlucassantos/vnc-domains/src/domains/deputy"
	"github.com/devlucassantos/vnc-domains/src/domains/eventsituation"
	"github.com/devlucassantos/vnc-domains/src/domains/eventtype"
	"github.com/devlucassantos/vnc-domains/src/domains/externalauthor"
	"github.com/devlucassantos/vnc-domains/src/domains/legislativebody"
	"github.com/devlucassantos/vnc-domains/src/domains/party"
	"github.com/devlucassantos/vnc-domains/src/domains/propositiontype"
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

func (instance Resources) GetResources() ([]articletype.ArticleType, []propositiontype.PropositionType, []party.Party,
	[]deputy.Deputy, []externalauthor.ExternalAuthor, []legislativebody.LegislativeBody, []eventtype.EventType,
	[]eventsituation.EventSituation, error) {
	articleTypes, err := instance.repository.GetArticleTypes()
	if err != nil {
		log.Errorf("Error retrieving article types data from the database: %s", err.Error())
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	propositionTypes, err := instance.repository.GetPropositionTypes()
	if err != nil {
		log.Errorf("Error retrieving proposition types data from the database: %s", err.Error())
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	parties, err := instance.repository.GetParties()
	if err != nil {
		log.Errorf("Error retrieving parties data from the database: %s", err.Error())
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	deputies, err := instance.repository.GetDeputies()
	if err != nil {
		log.Errorf("Error retrieving deputies data from the database: %s", err.Error())
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	externalAuthors, err := instance.repository.GetExternalAuthors()
	if err != nil {
		log.Errorf("Error retrieving external authors data from the database: %s", err.Error())
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	legislativeBodies, err := instance.repository.GetLegislativeBodies()
	if err != nil {
		log.Errorf("Error retrieving legislative bodies data from the database: %s", err.Error())
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	eventTypes, err := instance.repository.GetEventTypes()
	if err != nil {
		log.Errorf("Error retrieving event types data from the database: %s", err.Error())
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	eventSituations, err := instance.repository.GetEventSituations()
	if err != nil {
		log.Errorf("Error retrieving event situations data from the database: %s", err.Error())
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	return articleTypes, propositionTypes, parties, deputies, externalAuthors, legislativeBodies, eventTypes,
		eventSituations, nil
}

func (instance Resources) GetArticleTypes() ([]articletype.ArticleType, error) {
	articleTypes, err := instance.repository.GetArticleTypes()
	if err != nil {
		log.Errorf("Error retrieving article types data from the database: %s", err.Error())
		return nil, err
	}

	return articleTypes, nil
}

func (instance Resources) GetPropositionTypes() ([]propositiontype.PropositionType, error) {
	propositionTypes, err := instance.repository.GetPropositionTypes()
	if err != nil {
		log.Errorf("Error retrieving proposition types data from the database: %s", err.Error())
		return nil, err
	}

	return propositionTypes, nil
}

func (instance Resources) GetParties() ([]party.Party, error) {
	parties, err := instance.repository.GetParties()
	if err != nil {
		log.Errorf("Error retrieving parties data from the database: %s", err.Error())
		return nil, err
	}

	return parties, nil
}

func (instance Resources) GetDeputies() ([]deputy.Deputy, error) {
	deputies, err := instance.repository.GetDeputies()
	if err != nil {
		log.Errorf("Error retrieving deputies data from the database: %s", err.Error())
		return nil, err
	}

	return deputies, nil
}

func (instance Resources) GetExternalAuthors() ([]externalauthor.ExternalAuthor, error) {
	externalAuthors, err := instance.repository.GetExternalAuthors()
	if err != nil {
		log.Errorf("Error retrieving external authors data from the database: %s", err.Error())
		return nil, err
	}

	return externalAuthors, nil
}

func (instance Resources) GetLegislativeBodies() ([]legislativebody.LegislativeBody, error) {
	legislativeBodies, err := instance.repository.GetLegislativeBodies()
	if err != nil {
		log.Errorf("Error retrieving legislative bodies data from the database: %s", err.Error())
		return nil, err
	}

	return legislativeBodies, nil
}

func (instance Resources) GetEventTypes() ([]eventtype.EventType, error) {
	eventTypes, err := instance.repository.GetEventTypes()
	if err != nil {
		log.Errorf("Error retrieving event types data from the database: %s", err.Error())
		return nil, err
	}

	return eventTypes, nil
}

func (instance Resources) GetEventSituations() ([]eventsituation.EventSituation, error) {
	eventSituations, err := instance.repository.GetEventSituations()
	if err != nil {
		log.Errorf("Error retrieving event situations data from the database: %s", err.Error())
		return nil, err
	}

	return eventSituations, nil
}
