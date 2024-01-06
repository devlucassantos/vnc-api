package postgres

import (
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"vnc-read-api/core/domains/deputy"
	"vnc-read-api/core/domains/organization"
	"vnc-read-api/core/domains/party"
	"vnc-read-api/core/domains/proposition"
	"vnc-read-api/infra/dto"
	"vnc-read-api/infra/postgres/queries"
)

type Proposition struct {
	connectionManager ConnectionManagerInterface
}

func NewPropositionRepository(connectionManager ConnectionManagerInterface) *Proposition {
	return &Proposition{
		connectionManager: connectionManager,
	}
}

func (instance Proposition) GetPropositionById(id uuid.UUID) (*proposition.Proposition, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		return nil, err
	}
	defer instance.connectionManager.endConnection(postgresConnection)

	var propositionData dto.Proposition
	err = postgresConnection.Get(&propositionData, queries.Proposition().Select().ById(), id)
	if err != nil {
		log.Errorf("Erro ao obter os dados da proposição %s no banco de dados: %s", id, err.Error())
		return nil, err
	}

	var propositionAuthors []dto.Author
	err = postgresConnection.Select(&propositionAuthors, queries.PropositionAuthor().Select().ByPropositionId(), id)
	if err != nil {
		log.Errorf("Erro ao obter os dados dos autores da proposição %s no banco de dados: %s", id, err.Error())
		return nil, err
	}

	var deputies []deputy.Deputy
	var organizations []organization.Organization
	for _, author := range propositionAuthors {
		if author.Deputy.Id != uuid.Nil {
			currentParty, err := party.NewBuilder().
				Id(author.Deputy.Party.Id).
				Code(author.Deputy.Party.Code).
				Name(author.Deputy.Party.Name).
				Acronym(author.Deputy.Party.Acronym).
				ImageUrl(author.Deputy.Party.ImageUrl).
				Active(author.Deputy.Party.Active).
				CreatedAt(author.Deputy.Party.CreatedAt).
				UpdatedAt(author.Deputy.Party.UpdatedAt).
				Build()
			if err != nil {
				log.Errorf("Erro construindo a estrutura de dados do partido atual %s do(a) deputado(a) %s para a proposição %s: %s",
					author.Deputy.Party.Id, author.Deputy.Id, id, err.Error())
			}

			partyInTheProposition, err := party.NewBuilder().
				Id(author.Deputy.PartyInTheProposition.Id).
				Code(author.Deputy.PartyInTheProposition.Code).
				Name(author.Deputy.PartyInTheProposition.Name).
				Acronym(author.Deputy.PartyInTheProposition.Acronym).
				ImageUrl(author.Deputy.PartyInTheProposition.ImageUrl).
				Active(author.Deputy.PartyInTheProposition.Active).
				CreatedAt(author.Deputy.PartyInTheProposition.CreatedAt).
				UpdatedAt(author.Deputy.PartyInTheProposition.UpdatedAt).
				Build()
			if err != nil {
				log.Errorf("Erro construindo a estrutura de dados do partido na proposição %s do(a) deputado(a) "+
					"%s para a proposição %s: %s", author.Deputy.Party.Id, author.Deputy.Id, id, err.Error())
			}

			deputyDomain, err := deputy.NewBuilder().
				Id(author.Deputy.Id).
				Code(author.Deputy.Code).
				Cpf(author.Deputy.Cpf).
				Name(author.Deputy.Name).
				ElectoralName(author.Deputy.ElectoralName).
				ImageUrl(author.Deputy.ImageUrl).
				Party(*currentParty).
				PartyInTheProposition(*partyInTheProposition).
				Active(author.Deputy.Active).
				CreatedAt(author.Deputy.CreatedAt).
				UpdatedAt(author.Deputy.UpdatedAt).
				Build()
			if err != nil {
				log.Errorf("Erro construindo a estrutura de dados do(a) deputado(a) %s para a proposição %s: %s",
					author.Deputy.Id, id, err.Error())
				continue
			}

			deputies = append(deputies, *deputyDomain)
		} else if author.Organization.Id != uuid.Nil {
			organizationBuilder := organization.NewBuilder().
				Id(author.Organization.Id)

			if author.Organization.Code > 0 {
				organizationBuilder.Code(author.Organization.Code)
			}

			organizationDomain, err := organizationBuilder.
				Name(author.Organization.Name).
				Acronym(author.Organization.Acronym).
				Nickname(author.Organization.Nickname).
				Type(author.Organization.Type).
				Active(author.Organization.Active).
				CreatedAt(author.Organization.CreatedAt).
				UpdatedAt(author.Organization.UpdatedAt).
				Build()
			if err != nil {
				log.Errorf("Erro construindo a estrutura de dados da organização %s para a proposição %s: %s",
					author.Organization.Id, id, err.Error())
				continue
			}

			organizations = append(organizations, *organizationDomain)
		}
	}

	propositionDomain, err := proposition.NewBuilder().
		Id(propositionData.Id).
		Code(propositionData.Code).
		OriginalTextUrl(propositionData.OriginalTextUrl).
		Title(propositionData.Title).
		Content(propositionData.Content).
		SubmittedAt(propositionData.SubmittedAt).
		Deputies(deputies).
		Organizations(organizations).
		Active(propositionData.Active).
		CreatedAt(propositionData.CreatedAt).
		UpdatedAt(propositionData.UpdatedAt).
		Build()
	if err != nil {
		log.Errorf("Erro construindo a estrutura de dados da proposição %s: %s", id, err.Error())
		return nil, err
	}

	_, err = postgresConnection.Exec(queries.NewsView().Insert(), propositionData.NewsId)
	if err != nil {
		log.Errorf("Erro ao registrar a visualização da proposição %s: %s", propositionDomain.Id, err.Error())
	}

	return propositionDomain, nil
}
