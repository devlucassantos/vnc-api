package postgres

import (
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"vnc-read-api/core/domains/newsletter"
	"vnc-read-api/core/domains/proposition"
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

	var newsletterData dto.Newsletter
	err = postgresConnection.Get(&newsletterData, queries.Newsletter().Select().ById(), id)
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
		Id(newsletterData.Id).
		Title(newsletterData.Title).
		Content(newsletterData.Content).
		ReferenceDate(newsletterData.ReferenceDate).
		Propositions(propositions).
		Active(newsletterData.Active).
		CreatedAt(newsletterData.CreatedAt).
		UpdatedAt(newsletterData.UpdatedAt).
		Build()
	if err != nil {
		log.Errorf("Erro construindo a estrutura de dados do boletim %s: %s", id, err.Error())
		return nil, err
	}

	_, err = postgresConnection.Exec(queries.NewsView().Insert(), newsletterData.NewsId)
	if err != nil {
		log.Errorf("Erro ao registrar a visualização do boletim %s: %s", newsletterData.Id, err.Error())
	}

	return newsletterDomain, nil
}
