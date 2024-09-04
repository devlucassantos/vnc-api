package postgres

import (
	"database/sql"
	"github.com/devlucassantos/vnc-domains/src/domains/article"
	"github.com/devlucassantos/vnc-domains/src/domains/articletype"
	"github.com/devlucassantos/vnc-domains/src/domains/deputy"
	"github.com/devlucassantos/vnc-domains/src/domains/external"
	"github.com/devlucassantos/vnc-domains/src/domains/party"
	"github.com/devlucassantos/vnc-domains/src/domains/proposition"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"vnc-api/infra/dto"
	"vnc-api/infra/postgres/queries"
)

type Proposition struct {
	connectionManager connectionManagerInterface
}

func NewPropositionRepository(connectionManager connectionManagerInterface) *Proposition {
	return &Proposition{
		connectionManager: connectionManager,
	}
}

func (instance Proposition) GetPropositionByArticleId(articleId uuid.UUID, userId uuid.UUID) (*proposition.Proposition, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var propositionArticle dto.Article
	err = postgresConnection.Get(&propositionArticle, queries.Proposition().Select().ByArticleId(), articleId)
	if err != nil {
		log.Errorf("Erro ao obter os dados da proposição pela matéria %s no banco de dados: %s", articleId, err.Error())
		return nil, err
	}

	var propositionAuthors []dto.Author
	err = postgresConnection.Select(&propositionAuthors, queries.PropositionAuthor().Select().ByPropositionId(),
		propositionArticle.Proposition.Id)
	if err != nil {
		log.Errorf("Erro ao obter os dados dos autores da proposição %s no banco de dados: %s",
			propositionArticle.Proposition.Id, err.Error())
		return nil, err
	}

	var deputies []deputy.Deputy
	var externalAuthors []external.ExternalAuthor
	for _, author := range propositionAuthors {
		if author.Deputy.Id != uuid.Nil {
			currentParty, err := party.NewBuilder().
				Id(author.Deputy.Party.Id).
				Name(author.Deputy.Party.Name).
				Acronym(author.Deputy.Party.Acronym).
				ImageUrl(author.Deputy.Party.ImageUrl).
				CreatedAt(author.Deputy.Party.CreatedAt).
				UpdatedAt(author.Deputy.Party.UpdatedAt).
				Build()
			if err != nil {
				log.Errorf("Erro ao validar os dados do partido atual %s do(a) deputado(a) %s para a proposição %s: %s",
					author.Deputy.Party.Id, author.Deputy.Id, propositionArticle.Proposition.Id, err.Error())
			}

			partyInTheProposition, err := party.NewBuilder().
				Id(author.Deputy.PartyInTheProposition.Id).
				Name(author.Deputy.PartyInTheProposition.Name).
				Acronym(author.Deputy.PartyInTheProposition.Acronym).
				ImageUrl(author.Deputy.PartyInTheProposition.ImageUrl).
				CreatedAt(author.Deputy.PartyInTheProposition.CreatedAt).
				UpdatedAt(author.Deputy.PartyInTheProposition.UpdatedAt).
				Build()
			if err != nil {
				log.Errorf("Erro ao validar os dados do partido %s durante a proposição pelo(a) deputado(a) %s "+
					"para a proposição %s: %s", author.Deputy.Party.Id, author.Deputy.Id, propositionArticle.Proposition.Id, err.Error())
			}

			deputyDomain, err := deputy.NewBuilder().
				Id(author.Deputy.Id).
				Name(author.Deputy.Name).
				ElectoralName(author.Deputy.ElectoralName).
				ImageUrl(author.Deputy.ImageUrl).
				Party(*currentParty).
				PartyInTheProposition(*partyInTheProposition).
				CreatedAt(author.Deputy.CreatedAt).
				UpdatedAt(author.Deputy.UpdatedAt).
				Build()
			if err != nil {
				log.Errorf("Erro ao validar os dados do(a) deputado(a) %s para a proposição %s: %s",
					author.Deputy.Id, propositionArticle.Proposition.Id, err.Error())
				continue
			}

			deputies = append(deputies, *deputyDomain)
		} else if author.ExternalAuthor.Id != uuid.Nil {
			externalAuthor, err := external.NewBuilder().
				Id(author.ExternalAuthor.Id).
				Name(author.ExternalAuthor.Name).
				Type(author.ExternalAuthor.Type).
				CreatedAt(author.ExternalAuthor.CreatedAt).
				UpdatedAt(author.ExternalAuthor.UpdatedAt).
				Build()
			if err != nil {
				log.Errorf("Erro construindo a estrutura de dados do autor externo %s para a proposição %s: %s",
					author.ExternalAuthor.Id, propositionArticle.Proposition.Id, err.Error())
				continue
			}

			externalAuthors = append(externalAuthors, *externalAuthor)
		}
	}

	var userArticle dto.UserArticle
	if userId != uuid.Nil && propositionArticle.Id != uuid.Nil {
		numberOfArticles := 1
		err = postgresConnection.Get(&userArticle, queries.UserArticle().Select().RatingsAndArticlesSavedForLaterViewing(
			numberOfArticles), userId, propositionArticle.Id)
		if err != nil && err != sql.ErrNoRows {
			log.Errorf("Erro as buscar no banco de dados as informações da matéria %s que o usuário %s avaliou e/ou "+
				"salvou para ver depois: %s", articleId, userId, err.Error())
			return nil, err
		}
	}

	articleBuilder := article.NewBuilder()

	if userArticle.Article != nil {
		articleBuilder.UserRating(userArticle.Rating).ViewLater(userArticle.ViewLater)
	}

	articleType, err := articletype.NewBuilder().
		Id(propositionArticle.ArticleType.Id).
		Description(propositionArticle.ArticleType.Description).
		Color(propositionArticle.ArticleType.Color).
		SortOrder(propositionArticle.ArticleType.SortOrder).
		CreatedAt(propositionArticle.ArticleType.CreatedAt).
		UpdatedAt(propositionArticle.ArticleType.UpdatedAt).
		Build()
	if err != nil {
		log.Errorf("Erro ao validar os dados do tipo da matéria %s: %s", articleId, err.Error())
		return nil, err
	}

	articleDomain, err := articleBuilder.
		Id(propositionArticle.Id).
		AverageRating(propositionArticle.AverageRating).
		NumberOfRatings(propositionArticle.NumberOfRatings).
		Type(*articleType).
		ReferenceDateTime(propositionArticle.ReferenceDateTime).
		CreatedAt(propositionArticle.CreatedAt).
		UpdatedAt(propositionArticle.UpdatedAt).
		Build()
	if err != nil {
		log.Errorf("Erro ao validar os dados da matéria %s da proposição %s: %s", propositionArticle.Id,
			propositionArticle.Proposition.Id, err.Error())
		return nil, err
	}

	propositionDomain, err := proposition.NewBuilder().
		Id(propositionArticle.Proposition.Id).
		OriginalTextUrl(propositionArticle.Proposition.OriginalTextUrl).
		Title(propositionArticle.Proposition.Title).
		Content(propositionArticle.Proposition.Content).
		SubmittedAt(propositionArticle.Proposition.SubmittedAt).
		ImageUrl(propositionArticle.Proposition.ImageUrl).
		Deputies(deputies).
		ExternalAuthors(externalAuthors).
		Article(*articleDomain).
		CreatedAt(propositionArticle.Proposition.CreatedAt).
		UpdatedAt(propositionArticle.Proposition.UpdatedAt).
		Build()
	if err != nil {
		log.Errorf("Erro construindo a estrutura de dados da proposição %s: %s", articleId, err.Error())
		return nil, err
	}

	var userIdPointer *uuid.UUID
	if userId != uuid.Nil {
		userIdPointer = &userId
	}

	_, err = postgresConnection.Exec(queries.ArticleView().Insert(), propositionArticle.Id, userIdPointer)
	if err != nil {
		log.Errorf("Erro ao registrar a visualização da proposição %s: %s", articleId, err.Error())
	}

	return propositionDomain, nil
}
