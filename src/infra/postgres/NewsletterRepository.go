package postgres

import (
	"github.com/devlucassantos/vnc-domains/src/domains/newsletter"
	"github.com/devlucassantos/vnc-domains/src/domains/proposition"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"vnc-read-api/infra/dto"
	"vnc-read-api/infra/postgres/queries"
)

type Newsletter struct {
	connectionManager ConnectionManagerInterface
}

func NewNewsletterRepository(connectionManager ConnectionManagerInterface) *Newsletter {
	return &Newsletter{
		connectionManager: connectionManager,
	}
}

func (instance Newsletter) GetNewsletterById(id uuid.UUID) (*newsletter.Newsletter, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		return nil, err
	}
	defer instance.connectionManager.endConnection(postgresConnection)

	var newsData dto.News
	err = postgresConnection.Get(&newsData, queries.Newsletter().Select().ById(), id)
	if err != nil {
		log.Errorf("Erro ao obter os dados do boletim %s no banco de dados: %s", id, err.Error())
		return nil, err
	}

	var newsletterPropositions []dto.Proposition
	err = postgresConnection.Select(&newsletterPropositions, queries.NewsletterProposition().Select().ByNewsletterId(), id)
	if err != nil {
		log.Errorf("Erro ao obter os dados das proposições do boletim %s no banco de dados: %s", id, err.Error())
		return nil, err
	}

	var propositions []proposition.Proposition
	for _, propositionData := range newsletterPropositions {
		propositionDomain, err := proposition.NewBuilder().
			Id(propositionData.Id).
			Code(propositionData.Code).
			OriginalTextUrl(propositionData.OriginalTextUrl).
			Title(propositionData.Title).
			Content(propositionData.Content).
			SubmittedAt(propositionData.SubmittedAt).
			Active(propositionData.Active).
			CreatedAt(propositionData.CreatedAt).
			UpdatedAt(propositionData.UpdatedAt).
			Build()
		if err != nil {
			log.Errorf("Erro construindo a estrutura de dados da proposição %s: %s", id, err.Error())
			continue
		}

		propositions = append(propositions, *propositionDomain)
	}

	newsletterDomain, err := newsletter.NewBuilder().
		Id(newsData.Newsletter.Id).
		Title(newsData.Newsletter.Title).
		Content(newsData.Newsletter.Content).
		ReferenceDate(newsData.Newsletter.ReferenceDate).
		Propositions(propositions).
		Active(newsData.Newsletter.Active).
		CreatedAt(newsData.Newsletter.CreatedAt).
		UpdatedAt(newsData.Newsletter.UpdatedAt).
		Build()
	if err != nil {
		log.Errorf("Erro construindo a estrutura de dados do boletim %s: %s", id, err.Error())
		return nil, err
	}

	_, err = postgresConnection.Exec(queries.NewsView().Insert(), newsData.Id)
	if err != nil {
		log.Errorf("Erro ao registrar a visualização do boletim %s: %s", newsData.Id, err.Error())
	}

	return newsletterDomain, nil
}

func (instance Newsletter) GetNewsletterByPropositionId(propositionId uuid.UUID) (*newsletter.Newsletter, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		return nil, err
	}
	defer instance.connectionManager.endConnection(postgresConnection)

	var newsletterData dto.Newsletter
	err = postgresConnection.Get(&newsletterData, queries.Newsletter().Select().ByPropositionId(), propositionId)
	if err != nil {
		log.Errorf("Erro ao obter os dados do boletim relacionado a proposição %s no banco de dados: %s", propositionId, err.Error())
		return nil, err
	}

	newsletterDomain, err := newsletter.NewBuilder().
		Id(newsletterData.Id).
		Title(newsletterData.Title).
		Content(newsletterData.Content).
		ReferenceDate(newsletterData.ReferenceDate).
		Active(newsletterData.Active).
		CreatedAt(newsletterData.CreatedAt).
		UpdatedAt(newsletterData.UpdatedAt).
		Build()
	if err != nil {
		log.Errorf("Erro construindo a estrutura de dados do boletim relacionado a proposição %s: %s", propositionId, err.Error())
		return nil, err
	}

	return newsletterDomain, nil
}
