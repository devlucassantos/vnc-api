package postgres

import (
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"vnc-read-api/api/endpoints/dto/filter"
	"vnc-read-api/core/domains/deputy"
	"vnc-read-api/core/domains/keyword"
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

func (instance Proposition) GetPropositions(filter filter.PropositionFilter) ([]proposition.Proposition, int, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		return nil, 0, err
	}
	defer instance.connectionManager.endConnection(postgresConnection)

	var propositions []dto.Proposition
	err = postgresConnection.Select(&propositions, queries.Proposition().Select().All(),
		filter.PaginationFilter.CalculateOffset(), filter.PaginationFilter.GetItemsPerPage())
	if err != nil {
		log.Error("Erro ao obter os dados da proposições no banco de dados: ", err.Error())
		return nil, 0, err
	}

	var propositionData []proposition.Proposition
	for _, propositionDetails := range propositions {
		var propositionKeywords []dto.Keyword
		err = postgresConnection.Select(&propositionKeywords, queries.PropositionKeyword().Select().ByPropositionId(),
			propositionDetails.Id)
		if err != nil {
			log.Errorf("Erro ao obter os dados das palavras-chaves da proposição %s no banco de dados: %s",
				propositionDetails.Id, err.Error())
			return nil, 0, err
		}

		var propositionAuthors []dto.Author
		err = postgresConnection.Select(&propositionAuthors, queries.PropositionAuthor().Select().ByPropositionId(),
			propositionDetails.Id)
		if err != nil {
			log.Errorf("Erro ao obter os dados dos autores da proposição %s no banco de dados: %s",
				propositionDetails.Id, err.Error())
			return nil, 0, err
		}

		var keywords []keyword.Keyword
		for _, keywordData := range propositionKeywords {
			keywordDomain, err := keyword.NewBuilder().
				Id(keywordData.Id).
				Keyword(keywordData.Keyword).
				Active(keywordData.Active).
				CreatedAt(keywordData.CreatedAt).
				UpdatedAt(keywordData.UpdatedAt).
				Build()
			if err != nil {
				log.Errorf("Erro construindo a estrutura de dados da palavra-chave %s para a proposição %s: %s",
					keywordData.Keyword, propositionDetails.Id, err.Error())
				return nil, 0, err
			}
			keywords = append(keywords, *keywordDomain)
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
					log.Errorf("Erro construindo a estrutura de dados do partido %s do(a) deputado(a) %s para a proposição %s: %s",
						author.Deputy.Party.Id, author.Deputy.Id, propositionDetails.Id, err.Error())
					return nil, 0, err
				}

				deputyDomain, err := deputy.NewBuilder().
					Id(author.Deputy.Id).
					Code(author.Deputy.Code).
					Cpf(author.Deputy.Cpf).
					Name(author.Deputy.Name).
					ElectoralName(author.Deputy.ElectoralName).
					ImageUrl(author.Deputy.ImageUrl).
					CurrentParty(*currentParty).
					Active(author.Deputy.Active).
					CreatedAt(author.Deputy.CreatedAt).
					UpdatedAt(author.Deputy.UpdatedAt).
					Build()
				if err != nil {
					log.Errorf("Erro construindo a estrutura de dados do(a) deputado(a) %s para a proposição %s: %s",
						author.Deputy.Id, propositionDetails.Id, err.Error())
				}

				deputies = append(deputies, *deputyDomain)
			} else if author.Organization.Id != uuid.Nil {
				organizationDomain, err := organization.NewBuilder().
					Id(author.Organization.Id).
					Code(author.Organization.Code).
					Name(author.Organization.Name).
					Acronym(author.Organization.Acronym).
					Nickname(author.Organization.Nickname).
					Active(author.Organization.Active).
					CreatedAt(author.Organization.CreatedAt).
					UpdatedAt(author.Organization.UpdatedAt).
					Build()
				if err != nil {
					log.Errorf("Erro construindo a estrutura de dados da organização %s para a proposição %s: %s",
						author.Organization.Id, propositionDetails.Id, err.Error())
				}

				organizations = append(organizations, *organizationDomain)
			}
		}

		propositionDomain, err := proposition.NewBuilder().
			Id(propositionDetails.Id).
			Code(propositionDetails.Code).
			OriginalTextUrl(propositionDetails.OriginalTextUrl).
			Title(propositionDetails.Title).
			Summary(propositionDetails.Summary).
			SubmittedAt(propositionDetails.SubmittedAt).
			Deputies(deputies).
			Organizations(organizations).
			Keywords(keywords).
			Active(propositionDetails.Active).
			CreatedAt(propositionDetails.CreatedAt).
			UpdatedAt(propositionDetails.UpdatedAt).
			Build()
		if err != nil {
			log.Errorf("Erro construindo a estrutura de dados da proposição %s: %s", propositionDetails, err.Error())
			return nil, 0, err
		}

		propositionData = append(propositionData, *propositionDomain)
	}

	var totalNumberOfPropositions int
	err = postgresConnection.Get(&totalNumberOfPropositions, queries.Proposition().Select().TotalNumber())
	if err != nil {
		log.Error("Erro ao obter a quantidade total de proposições no banco de dados: ", err.Error())
		return nil, 0, err
	}

	return propositionData, totalNumberOfPropositions, nil
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

	var propositionKeywords []dto.Keyword
	err = postgresConnection.Select(&propositionKeywords, queries.PropositionKeyword().Select().ByPropositionId(), id)
	if err != nil {
		log.Errorf("Erro ao obter os dados das palavras-chaves da proposição %s no banco de dados: %s", id, err.Error())
		return nil, err
	}

	var propositionAuthors []dto.Author
	err = postgresConnection.Select(&propositionAuthors, queries.PropositionAuthor().Select().ByPropositionId(), id)
	if err != nil {
		log.Errorf("Erro ao obter os dados dos autores da proposição %s no banco de dados: %s", id, err.Error())
		return nil, err
	}

	var keywords []keyword.Keyword
	for _, keywordData := range propositionKeywords {
		keywordDomain, err := keyword.NewBuilder().
			Id(keywordData.Id).
			Keyword(keywordData.Keyword).
			Active(keywordData.Active).
			CreatedAt(keywordData.CreatedAt).
			UpdatedAt(keywordData.UpdatedAt).
			Build()
		if err != nil {
			log.Errorf("Erro construindo a estrutura de dados da palavra-chave %s para a proposição %s: %s",
				keywordData.Keyword, id, err.Error())
			return nil, err
		}
		keywords = append(keywords, *keywordDomain)
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
				log.Errorf("Erro construindo a estrutura de dados do partido %s do(a) deputado(a) %s para a proposição %s: %s",
					author.Deputy.Party.Id, author.Deputy.Id, id, err.Error())
				return nil, err
			}

			deputyDomain, err := deputy.NewBuilder().
				Id(author.Deputy.Id).
				Code(author.Deputy.Code).
				Cpf(author.Deputy.Cpf).
				Name(author.Deputy.Name).
				ElectoralName(author.Deputy.ElectoralName).
				ImageUrl(author.Deputy.ImageUrl).
				CurrentParty(*currentParty).
				Active(author.Deputy.Active).
				CreatedAt(author.Deputy.CreatedAt).
				UpdatedAt(author.Deputy.UpdatedAt).
				Build()
			if err != nil {
				log.Errorf("Erro construindo a estrutura de dados do(a) deputado(a) %s para a proposição %s: %s",
					author.Deputy.Id, id, err.Error())
			}

			deputies = append(deputies, *deputyDomain)
		} else if author.Organization.Id != uuid.Nil {
			organizationDomain, err := organization.NewBuilder().
				Id(author.Organization.Id).
				Code(author.Organization.Code).
				Name(author.Organization.Name).
				Acronym(author.Organization.Acronym).
				Nickname(author.Organization.Nickname).
				Active(author.Organization.Active).
				CreatedAt(author.Organization.CreatedAt).
				UpdatedAt(author.Organization.UpdatedAt).
				Build()
			if err != nil {
				log.Errorf("Erro construindo a estrutura de dados da organização %s para a proposição %s: %s",
					author.Organization.Id, id, err.Error())
			}

			organizations = append(organizations, *organizationDomain)
		}
	}

	propositionDomain, err := proposition.NewBuilder().
		Id(propositionData.Id).
		Code(propositionData.Code).
		OriginalTextUrl(propositionData.OriginalTextUrl).
		Title(propositionData.Title).
		Summary(propositionData.Summary).
		SubmittedAt(propositionData.SubmittedAt).
		Deputies(deputies).
		Organizations(organizations).
		Keywords(keywords).
		Active(propositionData.Active).
		CreatedAt(propositionData.CreatedAt).
		UpdatedAt(propositionData.UpdatedAt).
		Build()
	if err != nil {
		log.Errorf("Erro construindo a estrutura de dados da proposição %s: %s", id, err.Error())
		return nil, err
	}

	return propositionDomain, nil
}
