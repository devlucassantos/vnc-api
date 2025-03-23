package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/devlucassantos/vnc-domains/src/domains/agendaitemregime"
	"github.com/devlucassantos/vnc-domains/src/domains/article"
	"github.com/devlucassantos/vnc-domains/src/domains/articlesituation"
	"github.com/devlucassantos/vnc-domains/src/domains/articletype"
	"github.com/devlucassantos/vnc-domains/src/domains/deputy"
	"github.com/devlucassantos/vnc-domains/src/domains/event"
	"github.com/devlucassantos/vnc-domains/src/domains/eventagendaitem"
	"github.com/devlucassantos/vnc-domains/src/domains/eventsituation"
	"github.com/devlucassantos/vnc-domains/src/domains/eventtype"
	"github.com/devlucassantos/vnc-domains/src/domains/legislativebody"
	"github.com/devlucassantos/vnc-domains/src/domains/legislativebodytype"
	"github.com/devlucassantos/vnc-domains/src/domains/party"
	"github.com/devlucassantos/vnc-domains/src/domains/proposition"
	"github.com/devlucassantos/vnc-domains/src/domains/voting"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"vnc-api/infra/dto"
	"vnc-api/infra/postgres/queries"
)

type Event struct {
	connectionManager connectionManagerInterface
}

func NewEventRepository(connectionManager connectionManagerInterface) *Event {
	return &Event{
		connectionManager: connectionManager,
	}
}

func (instance Event) GetEventByArticleId(articleId uuid.UUID, userId uuid.UUID) (*event.Event, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("Error creating a connection to the Postgres database: ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var eventData dto.Article
	err = postgresConnection.Get(&eventData, queries.Event().Select().ByArticleId(), articleId)
	if err != nil {
		log.Errorf("Error retrieving event data for article %s from the database: %s", articleId, err.Error())
		return nil, err
	}

	eventType, err := eventtype.NewBuilder().
		Id(eventData.Event.EventType.Id).
		Description(eventData.Event.EventType.Description).
		Color(eventData.Event.EventType.Color).
		Build()
	if err != nil {
		log.Errorf("Error validating data for event type %s of event %s of article %s: %s",
			eventData.Event.EventType.Id, eventData.Event.Id, articleId, err.Error())
		return nil, err
	}

	eventSituation, err := eventsituation.NewBuilder().
		Id(eventData.Event.EventSituation.Id).
		Description(eventData.Event.EventSituation.Description).
		Color(eventData.Event.EventSituation.Color).
		Build()
	if err != nil {
		log.Errorf("Error validating data for event situation %s of event %s of article %s: %s",
			eventData.Event.EventSituation.Id, eventData.Event.Id, articleId, err.Error())
		return nil, err
	}

	var legislativeBodiesData []dto.LegislativeBody
	err = postgresConnection.Select(&legislativeBodiesData, queries.LegislativeBody().Select().
		LegislativeBodiesByEventId(), eventData.Event.Id)
	if err != nil {
		log.Errorf("Error retrieving legislative bodies data of event %s from the database: %s", articleId,
			err.Error())
		return nil, err
	}

	var legislativeBodies []legislativebody.LegislativeBody
	for _, legislativeBodyData := range legislativeBodiesData {
		legislativeBodyType, err := legislativebodytype.NewBuilder().
			Id(legislativeBodyData.LegislativeBodyType.Id).
			Description(legislativeBodyData.LegislativeBodyType.Description).
			Build()
		if err != nil {
			log.Errorf("Error validating data for the legislative body type %s of legislative body %s: %s",
				legislativeBodyData.LegislativeBodyType.Id, legislativeBodyData.Id, err.Error())
			return nil, err
		}

		legislativeBodyDomain, err := legislativebody.NewBuilder().
			Id(legislativeBodyData.Id).
			Name(legislativeBodyData.Name).
			Acronym(legislativeBodyData.Acronym).
			Type(*legislativeBodyType).
			Build()
		if err != nil {
			log.Errorf("Error validating data for the legislative body %s of event %s: %s",
				legislativeBodyData.Id, articleId, err.Error())
			return nil, err
		}
		legislativeBodies = append(legislativeBodies, *legislativeBodyDomain)
	}

	var requirementsData []dto.Article
	err = postgresConnection.Select(&requirementsData,
		queries.Article().Select().PropositionsOfTheRequirementsByEventId(), eventData.Event.Id)
	if err != nil {
		log.Errorf("Error retrieving data for articles of the requirements (propositions) of event %s "+
			"in the database: %s", articleId, err.Error())
		return nil, err
	}

	var agendaItemsData []dto.EventAgendaItem
	err = postgresConnection.Select(&agendaItemsData, queries.EventAgendaItem().Select().ByEventId(), eventData.Event.Id)
	if err != nil {
		log.Errorf("Error retrieving data for agenda items of event %s from the database: %s", articleId,
			err.Error())
		return nil, err
	}

	var agendaItemArticleIds []interface{}
	for _, agendaItemData := range agendaItemsData {
		agendaItemArticleIds = append(agendaItemArticleIds, agendaItemData.PropositionArticleId)

		if agendaItemData.RelatedPropositionArticleId != uuid.Nil {
			agendaItemArticleIds = append(agendaItemArticleIds, agendaItemData.RelatedPropositionArticleId)
		}

		if agendaItemData.VotingArticleId != uuid.Nil {
			agendaItemArticleIds = append(agendaItemArticleIds, agendaItemData.VotingArticleId)
		}
	}

	var agendaItemArticlesData []dto.Article
	if agendaItemArticleIds != nil {
		err = postgresConnection.Select(&agendaItemArticlesData, queries.Article().Select().In(len(agendaItemArticleIds)),
			agendaItemArticleIds...)
		if err != nil {
			log.Errorf("Error retrieving data for agenda item articles of event %s from the database: %s", articleId,
				err.Error())
			return nil, err
		}
	}

	agendaItemArticles := make(map[uuid.UUID]dto.Article)
	for _, agendaItemArticle := range agendaItemArticlesData {
		agendaItemArticles[agendaItemArticle.Id] = agendaItemArticle
	}

	userArticles := make(map[uuid.UUID]dto.UserArticle)
	if userId != uuid.Nil {
		eventArticles := append([]interface{}{}, eventData.Id)
		eventArticles = append(eventArticles, agendaItemArticleIds...)
		for _, requirementData := range requirementsData {
			eventArticles = append(eventArticles, requirementData.Id)
		}
		articleFilters := append([]interface{}{}, userId)
		articleFilters = append(articleFilters, eventArticles...)

		var userArticleData []dto.UserArticle
		err = postgresConnection.Select(&userArticleData, queries.Article().Select().
			RatingsAndArticlesSavedForLaterViewing(len(eventArticles)), articleFilters...)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			log.Errorf("Error retrieving articles that the user %s rated and/or saved for later viewing: %s",
				userId, err.Error())
			return nil, err
		}

		for _, userArticle := range userArticleData {
			userArticles[userArticle.Article.Id] = userArticle
		}
	}

	var requirements []proposition.Proposition
	for _, articleData := range requirementsData {
		articleType, err := articletype.NewBuilder().
			Id(articleData.ArticleType.Id).
			Description(articleData.ArticleType.Description).
			Codes(articleData.ArticleType.Codes).
			Color(articleData.ArticleType.Color).
			Build()
		if err != nil {
			log.Errorf("Error validating data for article type %s of article %s: %s", articleData.ArticleType.Id,
				articleData.Id, err.Error())
			return nil, err
		}

		articleSpecificType, err := articletype.NewBuilder().
			Id(articleData.Proposition.PropositionType.Id).
			Description(articleData.Proposition.PropositionType.Description).
			Color(articleData.Proposition.PropositionType.Color).
			Build()
		if err != nil {
			log.Errorf("Error validating data for proposition type %s of proposition %s of article %s: %s",
				articleData.Proposition.PropositionType.Id, articleData.Proposition.Id, articleData.Id, err.Error())
			return nil, err
		}

		articleBuilder := article.NewBuilder()

		if _, exists := userArticles[articleData.Id]; exists {
			articleBuilder.UserRating(userArticles[articleData.Id].Rating).
				ViewLater(userArticles[articleData.Id].ViewLater)
		}

		if articleData.Proposition.ImageUrl != "" {
			articleBuilder.MultimediaUrl(articleData.Proposition.ImageUrl).
				MultimediaDescription(articleData.Proposition.ImageDescription)
		}

		articleDomain, err := articleBuilder.
			Id(articleData.Id).
			Title(articleData.Proposition.Title).
			Content(articleData.Proposition.Content).
			AverageRating(articleData.AverageRating).
			NumberOfRatings(articleData.NumberOfRatings).
			Type(*articleType).
			SpecificType(*articleSpecificType).
			CreatedAt(articleData.CreatedAt).
			UpdatedAt(articleData.UpdatedAt).
			Build()
		if err != nil {
			log.Errorf("Error validating data for requirement article %s of event %s: %s", articleData.Id,
				articleId, err.Error())
			return nil, err
		}

		requirement, err := proposition.NewBuilder().
			Id(articleData.Proposition.Id).
			Article(*articleDomain).
			Build()
		if err != nil {
			log.Errorf("Error validating data for proposition %s of article %s of event %s (requirement): %s",
				articleData.Proposition.Id, articleData.Id, articleId, err)
			return nil, err
		}

		requirements = append(requirements, *requirement)
	}

	var agendaItems []eventagendaitem.EventAgendaItem
	for _, agendaItemData := range agendaItemsData {
		agendaItemRegime, err := agendaitemregime.NewBuilder().
			Id(agendaItemData.AgendaItemRegime.Id).
			Description(agendaItemData.AgendaItemRegime.Description).
			Build()
		if err != nil {
			log.Errorf("Error validating data for agenda item regime %s of agenda item %s of event %s: %s",
				agendaItemData.AgendaItemRegime.Id, agendaItemData.Id, articleId, err.Error())
			return nil, err
		}

		var rapporteur *deputy.Deputy
		if agendaItemData.Deputy.Id != uuid.Nil {
			currentParty, err := party.NewBuilder().
				Id(agendaItemData.Deputy.Party.Id).
				Name(agendaItemData.Deputy.Party.Name).
				Acronym(agendaItemData.Deputy.Party.Acronym).
				ImageUrl(agendaItemData.Deputy.Party.ImageUrl).
				ImageDescription(fmt.Sprintf("Logo do %s (%s)", agendaItemData.Deputy.Party.Name,
					agendaItemData.Deputy.Party.Acronym)).
				Build()
			if err != nil {
				log.Errorf("Error validating data for the current party %s of deputy %s for agenda item %s of "+
					"event %s: %s", agendaItemData.Deputy.Party.Id, agendaItemData.Deputy.Id, agendaItemData.Id,
					articleId, err.Error())
				return nil, err
			}

			previousParty, err := party.NewBuilder().
				Id(agendaItemData.Deputy.PreviousParty.Id).
				Name(agendaItemData.Deputy.PreviousParty.Name).
				Acronym(agendaItemData.Deputy.PreviousParty.Acronym).
				ImageUrl(agendaItemData.Deputy.PreviousParty.ImageUrl).
				ImageDescription(fmt.Sprintf("Logo do %s (%s)", agendaItemData.Deputy.PreviousParty.Name,
					agendaItemData.Deputy.PreviousParty.Acronym)).
				Build()
			if err != nil {
				log.Errorf("rror validating data for party %s of deputy %s when discussing agenda item %s of "+
					"event %s: %s", agendaItemData.Deputy.PreviousParty.Id, agendaItemData.Deputy.Id, agendaItemData.Id,
					articleId, err.Error())
				return nil, err
			}

			rapporteur, err = deputy.NewBuilder().
				Id(agendaItemData.Deputy.Id).
				Name(agendaItemData.Deputy.Name).
				ElectoralName(agendaItemData.Deputy.ElectoralName).
				ImageUrl(agendaItemData.Deputy.ImageUrl).
				ImageDescription(fmt.Sprintf("Foto do(a) deputado(a) federal %s (%s-%s)",
					agendaItemData.Deputy.Name, agendaItemData.Deputy.Party.Acronym,
					agendaItemData.Deputy.FederatedUnit)).
				Party(*currentParty).
				FederatedUnit(agendaItemData.Deputy.FederatedUnit).
				PreviousParty(*previousParty).
				PreviousFederatedUnit(agendaItemData.Deputy.PreviousFederatedUnit).
				Build()
			if err != nil {
				log.Errorf("Error validating data for deputy %s for agenda item %s of event %s: %s",
					agendaItemData.Deputy.Id, agendaItemData.Id, articleId, err.Error())
				return nil, err
			}
		}

		articleType, err := articletype.NewBuilder().
			Id(agendaItemArticles[agendaItemData.PropositionArticleId].ArticleType.Id).
			Description(agendaItemArticles[agendaItemData.PropositionArticleId].ArticleType.Description).
			Codes(agendaItemArticles[agendaItemData.PropositionArticleId].ArticleType.Codes).
			Color(agendaItemArticles[agendaItemData.PropositionArticleId].ArticleType.Color).
			Build()
		if err != nil {
			log.Errorf("Error validating data for article type %s of article %s: %s",
				agendaItemArticles[agendaItemData.PropositionArticleId].ArticleType.Id,
				agendaItemArticles[agendaItemData.PropositionArticleId].Id, err.Error())
			return nil, err
		}

		articleSpecificType, err := articletype.NewBuilder().
			Id(agendaItemArticles[agendaItemData.PropositionArticleId].Proposition.PropositionType.Id).
			Description(agendaItemArticles[agendaItemData.PropositionArticleId].Proposition.PropositionType.Description).
			Color(agendaItemArticles[agendaItemData.PropositionArticleId].Proposition.PropositionType.Color).
			Build()
		if err != nil {
			log.Errorf("Error validating data for proposition type %s of proposition %s of article %s: %s",
				agendaItemArticles[agendaItemData.PropositionArticleId].Proposition.PropositionType.Id,
				agendaItemArticles[agendaItemData.PropositionArticleId].Proposition.Id,
				agendaItemArticles[agendaItemData.PropositionArticleId].Id, err.Error())
			return nil, err
		}

		articleBuilder := article.NewBuilder()

		if _, exists := userArticles[agendaItemArticles[agendaItemData.PropositionArticleId].Id]; exists {
			articleBuilder.UserRating(userArticles[agendaItemArticles[agendaItemData.PropositionArticleId].Id].Rating).
				ViewLater(userArticles[agendaItemArticles[agendaItemData.PropositionArticleId].Id].ViewLater)
		}

		articleDomain, err := articleBuilder.
			Id(agendaItemArticles[agendaItemData.PropositionArticleId].Id).
			Title(agendaItemArticles[agendaItemData.PropositionArticleId].Proposition.Title).
			Content(agendaItemArticles[agendaItemData.PropositionArticleId].Proposition.Content).
			AverageRating(agendaItemArticles[agendaItemData.PropositionArticleId].AverageRating).
			NumberOfRatings(agendaItemArticles[agendaItemData.PropositionArticleId].NumberOfRatings).
			Type(*articleType).
			SpecificType(*articleSpecificType).
			CreatedAt(agendaItemArticles[agendaItemData.PropositionArticleId].CreatedAt).
			UpdatedAt(agendaItemArticles[agendaItemData.PropositionArticleId].UpdatedAt).
			Build()
		if err != nil {
			log.Errorf("Error validating data for article %s of proposition %s of agenda item %s of event %s: %s",
				agendaItemArticles[agendaItemData.PropositionArticleId].Id,
				agendaItemArticles[agendaItemData.PropositionArticleId].Proposition.Id, agendaItemData.Id, articleId,
				err.Error())
			return nil, err
		}

		agendaItemProposition, err := proposition.NewBuilder().
			Id(agendaItemArticles[agendaItemData.PropositionArticleId].Proposition.Id).
			Article(*articleDomain).
			Build()
		if err != nil {
			log.Errorf("Error validating data for proposition %s (Article ID: %s) of agenda item %s of event "+
				"%s: %s", agendaItemArticles[agendaItemData.PropositionArticleId].Proposition.Id,
				agendaItemArticles[agendaItemData.PropositionArticleId].Id, agendaItemData.Id, articleId, err.Error())
			return nil, err
		}

		var propositionRelatedToAgendaItem *proposition.Proposition
		if agendaItemData.RelatedPropositionArticleId != uuid.Nil {
			articleType, err = articletype.NewBuilder().
				Id(agendaItemArticles[agendaItemData.RelatedPropositionArticleId].ArticleType.Id).
				Description(agendaItemArticles[agendaItemData.RelatedPropositionArticleId].ArticleType.Description).
				Codes(agendaItemArticles[agendaItemData.RelatedPropositionArticleId].ArticleType.Codes).
				Color(agendaItemArticles[agendaItemData.RelatedPropositionArticleId].ArticleType.Color).
				Build()
			if err != nil {
				log.Errorf("Error validating data for article type %s of article %s: %s",
					agendaItemArticles[agendaItemData.RelatedPropositionArticleId].ArticleType.Id,
					agendaItemArticles[agendaItemData.RelatedPropositionArticleId].Id, err.Error())
				return nil, err
			}

			articleSpecificType, err = articletype.NewBuilder().
				Id(agendaItemArticles[agendaItemData.RelatedPropositionArticleId].Proposition.PropositionType.Id).
				Description(agendaItemArticles[agendaItemData.RelatedPropositionArticleId].Proposition.PropositionType.
					Description).
				Color(agendaItemArticles[agendaItemData.RelatedPropositionArticleId].Proposition.PropositionType.Color).
				Build()
			if err != nil {
				log.Errorf("Error validating data for proposition type %s of proposition %s of article %s: %s",
					agendaItemArticles[agendaItemData.RelatedPropositionArticleId].Proposition.PropositionType.Id,
					agendaItemArticles[agendaItemData.RelatedPropositionArticleId].Proposition.Id,
					agendaItemArticles[agendaItemData.RelatedPropositionArticleId].Id, err.Error())
				return nil, err
			}

			articleBuilder = article.NewBuilder()

			if _, exists := userArticles[agendaItemArticles[agendaItemData.RelatedPropositionArticleId].Id]; exists {
				articleBuilder.UserRating(userArticles[agendaItemArticles[agendaItemData.RelatedPropositionArticleId].
					Id].Rating).
					ViewLater(userArticles[agendaItemArticles[agendaItemData.RelatedPropositionArticleId].Id].ViewLater)
			}

			if agendaItemArticles[agendaItemData.RelatedPropositionArticleId].Proposition.ImageUrl != "" {
				articleBuilder.MultimediaUrl(agendaItemArticles[agendaItemData.RelatedPropositionArticleId].Proposition.
					ImageUrl).
					MultimediaDescription(agendaItemArticles[agendaItemData.RelatedPropositionArticleId].Proposition.
						ImageDescription)
			}

			articleDomain, err = articleBuilder.
				Id(agendaItemArticles[agendaItemData.RelatedPropositionArticleId].Id).
				Title(agendaItemArticles[agendaItemData.RelatedPropositionArticleId].Proposition.Title).
				Content(agendaItemArticles[agendaItemData.RelatedPropositionArticleId].Proposition.Content).
				AverageRating(agendaItemArticles[agendaItemData.RelatedPropositionArticleId].AverageRating).
				NumberOfRatings(agendaItemArticles[agendaItemData.RelatedPropositionArticleId].NumberOfRatings).
				Type(*articleType).
				SpecificType(*articleSpecificType).
				CreatedAt(agendaItemArticles[agendaItemData.RelatedPropositionArticleId].CreatedAt).
				UpdatedAt(agendaItemArticles[agendaItemData.RelatedPropositionArticleId].UpdatedAt).
				Build()
			if err != nil {
				log.Errorf("Error validating data for article %s of proposition %s related to agenda item %s of "+
					"event %s: %s", agendaItemArticles[agendaItemData.RelatedPropositionArticleId].Id,
					agendaItemArticles[agendaItemData.RelatedPropositionArticleId].Proposition.Id, agendaItemData.Id,
					articleId, err.Error())
				return nil, err
			}

			propositionRelatedToAgendaItem, err = proposition.NewBuilder().
				Id(agendaItemArticles[agendaItemData.RelatedPropositionArticleId].Proposition.Id).
				Article(*articleDomain).
				Build()
			if err != nil {
				log.Errorf("Error validating data for proposition %s (Article ID: %s) related to agenda item %s "+
					"of event %s: %s", agendaItemArticles[agendaItemData.RelatedPropositionArticleId].Proposition.Id,
					agendaItemArticles[agendaItemData.RelatedPropositionArticleId].Id, agendaItemData.Id, articleId,
					err.Error())
				return nil, err
			}
		}

		var agendaItemVoting *voting.Voting
		if agendaItemData.VotingArticleId != uuid.Nil {
			articleType, err = articletype.NewBuilder().
				Id(agendaItemArticles[agendaItemData.VotingArticleId].ArticleType.Id).
				Description(agendaItemArticles[agendaItemData.VotingArticleId].ArticleType.Description).
				Codes(agendaItemArticles[agendaItemData.VotingArticleId].ArticleType.Codes).
				Color(agendaItemArticles[agendaItemData.VotingArticleId].ArticleType.Color).
				Build()
			if err != nil {
				log.Errorf("Error validating data for article type %s of article %s: %s",
					agendaItemArticles[agendaItemData.VotingArticleId].ArticleType.Id,
					agendaItemArticles[agendaItemData.VotingArticleId].Id, err.Error())
				return nil, err
			}

			articleBuilder = article.NewBuilder()

			if _, exists := userArticles[agendaItemArticles[agendaItemData.VotingArticleId].Id]; exists {
				articleBuilder.UserRating(userArticles[agendaItemArticles[agendaItemData.VotingArticleId].Id].Rating).
					ViewLater(userArticles[agendaItemArticles[agendaItemData.VotingArticleId].Id].ViewLater)
			}

			articleSituation, err := articlesituation.NewBuilder().
				IsApproved(agendaItemArticles[agendaItemData.VotingArticleId].Voting.IsApproved).
				Build()
			if err != nil {
				log.Errorf("Error validating data for article situation of article %s: %s",
					agendaItemArticles[agendaItemData.VotingArticleId].Id, err.Error())
			}

			articleDomain, err = articleBuilder.
				Id(agendaItemArticles[agendaItemData.VotingArticleId].Id).
				Title(fmt.Sprint("Votação ", agendaItemArticles[agendaItemData.VotingArticleId].Voting.Code)).
				Content(agendaItemArticles[agendaItemData.VotingArticleId].Voting.Result).
				Situation(*articleSituation).
				AverageRating(agendaItemArticles[agendaItemData.VotingArticleId].AverageRating).
				NumberOfRatings(agendaItemArticles[agendaItemData.VotingArticleId].NumberOfRatings).
				Type(*articleType).
				CreatedAt(agendaItemArticles[agendaItemData.VotingArticleId].CreatedAt).
				UpdatedAt(agendaItemArticles[agendaItemData.VotingArticleId].UpdatedAt).
				Build()
			if err != nil {
				log.Errorf("Error validating data for article %s of voting %s of agenda item %s of event %s: %s",
					agendaItemArticles[agendaItemData.VotingArticleId].Id,
					agendaItemArticles[agendaItemData.VotingArticleId].Voting.Id, agendaItemData.Id, articleId,
					err.Error())
				return nil, err
			}

			agendaItemVoting, err = voting.NewBuilder().
				Id(agendaItemArticles[agendaItemData.VotingArticleId].Voting.Id).
				Article(*articleDomain).
				Build()
			if err != nil {
				log.Errorf("Error validating data for voting %s of article %s of agenda item %s of event %s: %s",
					agendaItemArticles[agendaItemData.VotingArticleId].Voting.Id,
					agendaItemArticles[agendaItemData.VotingArticleId].Id, agendaItemData.Id, articleId, err.Error())
				return nil, err
			}
		}

		agendaItemBuilder := eventagendaitem.NewBuilder()

		if agendaItemData.Situation != "" {
			agendaItemBuilder.Situation(agendaItemData.Situation)
		}

		if rapporteur != nil {
			agendaItemBuilder.Rapporteur(*rapporteur)
		}

		if propositionRelatedToAgendaItem != nil {
			agendaItemBuilder.RelatedProposition(*propositionRelatedToAgendaItem)
		}

		if agendaItemVoting != nil {
			agendaItemBuilder.Voting(*agendaItemVoting)
		}

		agendaItem, err := agendaItemBuilder.
			Id(agendaItemData.Id).
			Title(agendaItemData.Title).
			Topic(agendaItemData.Topic).
			Regime(*agendaItemRegime).
			Proposition(*agendaItemProposition).
			Build()
		if err != nil {
			log.Errorf("Error validating data for agenda item %s of event %s: %s", agendaItemData.Id, articleId,
				err.Error())
			return nil, err
		}

		agendaItems = append(agendaItems, *agendaItem)
	}

	articleBuilder := article.NewBuilder()

	if _, exists := userArticles[eventData.Id]; exists {
		articleBuilder.UserRating(userArticles[eventData.Id].Rating).
			ViewLater(userArticles[eventData.Id].ViewLater)
	}

	articleType, err := articletype.NewBuilder().
		Id(eventData.ArticleType.Id).
		Description(eventData.ArticleType.Description).
		Codes(eventData.ArticleType.Codes).
		Color(eventData.ArticleType.Color).
		Build()
	if err != nil {
		log.Errorf("Error validating data for article type %s of article %s: %s", eventData.ArticleType.Id,
			eventData.Id, err.Error())
		return nil, err
	}

	articleDomain, err := articleBuilder.
		Id(eventData.Id).
		AverageRating(eventData.AverageRating).
		NumberOfRatings(eventData.NumberOfRatings).
		Type(*articleType).
		CreatedAt(eventData.CreatedAt).
		UpdatedAt(eventData.UpdatedAt).
		Build()
	if err != nil {
		log.Errorf("Error validating data for article %s of event %s: %s", eventData.Id, eventData.Voting.Id,
			err.Error())
		return nil, err
	}

	eventBuilder := event.NewBuilder()

	if !eventData.Event.EndsAt.IsZero() {
		eventBuilder.EndsAt(eventData.Event.EndsAt)
	}

	if eventData.Event.VideoUrl != "" {
		eventBuilder.VideoUrl(eventData.Event.VideoUrl)
	}

	eventDomain, err := eventBuilder.
		Id(eventData.Event.Id).
		Code(eventData.Event.Code).
		Title(eventData.Event.Title).
		Description(eventData.Event.Description).
		StartsAt(eventData.Event.StartsAt).
		Location(eventData.Event.Location).
		IsInternal(eventData.Event.IsInternal).
		Type(*eventType).
		Situation(*eventSituation).
		LegislativeBodies(legislativeBodies).
		Requirements(requirements).
		AgendaItems(agendaItems).
		Article(*articleDomain).
		Build()
	if err != nil {
		log.Errorf("Error validating data for event %s of article %s: %s", eventData.Event.Id, articleId,
			err.Error())
		return nil, err
	}

	var userIdPointer *uuid.UUID
	if userId != uuid.Nil {
		userIdPointer = &userId
	}

	_, err = postgresConnection.Exec(queries.ArticleView().Insert(), articleId, userIdPointer)
	if err != nil {
		log.Errorf("Error registering the view for article %s: %s", articleId, err.Error())
	}

	return eventDomain, nil
}
