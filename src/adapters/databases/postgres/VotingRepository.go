package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/devlucassantos/vnc-domains/src/domains/article"
	"github.com/devlucassantos/vnc-domains/src/domains/articlesituation"
	"github.com/devlucassantos/vnc-domains/src/domains/articletype"
	"github.com/devlucassantos/vnc-domains/src/domains/legislativebody"
	"github.com/devlucassantos/vnc-domains/src/domains/legislativebodytype"
	"github.com/devlucassantos/vnc-domains/src/domains/proposition"
	"github.com/devlucassantos/vnc-domains/src/domains/voting"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"vnc-api/adapters/databases/dto"
	"vnc-api/adapters/databases/postgres/queries"
)

type Voting struct {
	connectionManager connectionManagerInterface
}

func NewVotingRepository(connectionManager connectionManagerInterface) *Voting {
	return &Voting{
		connectionManager: connectionManager,
	}
}

func (instance Voting) GetVotingByArticleId(articleId uuid.UUID, userId uuid.UUID) (*voting.Voting, error) {
	postgresConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("Error creating a connection to the Postgres database: ", err.Error())
		return nil, err
	}
	defer instance.connectionManager.closeConnection(postgresConnection)

	var votingArticle dto.Article
	err = postgresConnection.Get(&votingArticle, queries.Voting().Select().ByArticleId(), articleId)
	if err != nil {
		log.Errorf("Error retrieving voting data for article %s from the database: %s", articleId, err.Error())
		return nil, err
	}

	legislativeBodyType, err := legislativebodytype.NewBuilder().
		Id(votingArticle.Voting.LegislativeBody.LegislativeBodyType.Id).
		Description(votingArticle.Voting.LegislativeBody.LegislativeBodyType.Description).
		Build()
	if err != nil {
		log.Errorf("Error validating data for the legislative body type %s of legislative body %s of voting "+
			"%s: %s", votingArticle.Voting.LegislativeBody.LegislativeBodyType.Id, votingArticle.Voting.LegislativeBody.Id,
			articleId, err.Error())
		return nil, err
	}

	legislativeBody, err := legislativebody.NewBuilder().
		Id(votingArticle.Voting.LegislativeBody.Id).
		Name(votingArticle.Voting.LegislativeBody.Name).
		Acronym(votingArticle.Voting.LegislativeBody.Acronym).
		Type(*legislativeBodyType).
		Build()
	if err != nil {
		log.Errorf("Error validating data for legislative body %s of voting %s: %s",
			votingArticle.Voting.LegislativeBody.Id, articleId, err.Error())
		return nil, err
	}

	var mainPropositionData dto.Article
	if votingArticle.Voting.MainPropositionId != uuid.Nil {
		err = postgresConnection.Get(&mainPropositionData, queries.Article().Select().MainPropositionByVotingId(),
			votingArticle.Voting.Id)
		if err != nil {
			log.Errorf("Error retrieving data for main proposition of voting %s: %s", articleId, err.Error())
			return nil, err
		}
	}

	var articlesOfThePropositionsRelatedToVoting []dto.Article
	err = postgresConnection.Select(&articlesOfThePropositionsRelatedToVoting,
		queries.Article().Select().PropositionsRelatedByVotingId(), votingArticle.Voting.Id)
	if err != nil {
		log.Errorf("Error retrieving data for articles of the propositions related to voting %s from the "+
			"database: %s", articleId, err.Error())
		return nil, err
	}

	var articlesOfThePropositionsAffectedByVoting []dto.Article
	err = postgresConnection.Select(&articlesOfThePropositionsAffectedByVoting,
		queries.Article().Select().PropositionsAffectedByVotingId(), articleId)
	if err != nil {
		log.Errorf("Error retrieving data for articles of the propositions affected by voting %s from the "+
			"database: %s", articleId, err.Error())
		return nil, err
	}

	var relatedArticleIds []uuid.UUID
	err = postgresConnection.Select(&relatedArticleIds, queries.Article().Select().RelatedArticlesByVotingId(),
		votingArticle.Voting.Id)
	if err != nil {
		log.Errorf("Error retrieving articles IDs related to voting %s from the database: %s", articleId,
			err.Error())
		return nil, err
	}

	var relatedArticles []dto.Article
	if relatedArticleIds != nil {
		var relatedArticleIdsAsInterfaceSlice []interface{}
		for _, relatedArticleId := range relatedArticleIds {
			relatedArticleIdsAsInterfaceSlice = append(relatedArticleIdsAsInterfaceSlice, relatedArticleId)
		}

		err = postgresConnection.Select(&relatedArticles, queries.Article().Select().In(len(relatedArticleIds)),
			relatedArticleIdsAsInterfaceSlice...)
		if err != nil {
			log.Errorf("Error retrieving data for articles related to proposition %s from the database: %s",
				articleId, err.Error())
			return nil, err
		}
	}

	userArticles := make(map[uuid.UUID]dto.UserArticle)
	if userId != uuid.Nil {
		var votingArticles []interface{}
		votingArticles = append(votingArticles, votingArticle.Id)
		for _, relatedArticleId := range relatedArticleIds {
			votingArticles = append(votingArticles, relatedArticleId)
		}
		if mainPropositionData.Id != uuid.Nil {
			votingArticles = append(votingArticles, mainPropositionData.Id)
		}
		for _, articleData := range articlesOfThePropositionsRelatedToVoting {
			votingArticles = append(votingArticles, articleData.Id)
		}
		for _, articleData := range articlesOfThePropositionsAffectedByVoting {
			votingArticles = append(votingArticles, articleData.Id)
		}
		articleFilters := append([]interface{}{}, userId)
		articleFilters = append(articleFilters, votingArticles...)

		var userArticleData []dto.UserArticle
		err = postgresConnection.Select(&userArticleData, queries.Article().Select().
			RatingsAndArticlesSavedForLaterViewing(len(votingArticles)), articleFilters...)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			log.Errorf("Error retrieving articles that the user %s rated and/or saved for later viewing: %s",
				userId, err.Error())
			return nil, err
		}

		for _, userArticle := range userArticleData {
			userArticles[userArticle.Article.Id] = userArticle
		}
	}

	var mainProposition *proposition.Proposition
	if mainPropositionData.Id != uuid.Nil {
		articleType, err := articletype.NewBuilder().
			Id(mainPropositionData.ArticleType.Id).
			Description(mainPropositionData.ArticleType.Description).
			Codes(mainPropositionData.ArticleType.Codes).
			Color(mainPropositionData.ArticleType.Color).
			Build()
		if err != nil {
			log.Errorf("Error validating data for article type %s of article %s: %s",
				mainPropositionData.ArticleType.Id, mainPropositionData.Id, err.Error())
			return nil, err
		}

		articleSpecificType, err := articletype.NewBuilder().
			Id(mainPropositionData.Proposition.PropositionType.Id).
			Description(mainPropositionData.Proposition.PropositionType.Description).
			Color(mainPropositionData.Proposition.PropositionType.Color).
			Build()
		if err != nil {
			log.Errorf("Error validating data for proposition type %s of proposition %s of article %s: %s",
				mainPropositionData.Proposition.PropositionType.Id, mainPropositionData.Proposition.Id,
				mainPropositionData.Id, err.Error())
			return nil, err
		}

		articleBuilder := article.NewBuilder()

		if _, exists := userArticles[mainPropositionData.Id]; exists {
			articleBuilder.UserRating(userArticles[mainPropositionData.Id].Rating).
				ViewLater(userArticles[mainPropositionData.Id].ViewLater)
		}

		if mainPropositionData.Proposition.ImageUrl != "" {
			articleBuilder.MultimediaUrl(mainPropositionData.Proposition.ImageUrl).
				MultimediaDescription(mainPropositionData.Proposition.ImageDescription)
		}

		articleDomain, err := articleBuilder.
			Id(mainPropositionData.Id).
			Title(mainPropositionData.Proposition.Title).
			Content(mainPropositionData.Proposition.Content).
			AverageRating(mainPropositionData.AverageRating).
			NumberOfRatings(mainPropositionData.NumberOfRatings).
			Type(*articleType).
			SpecificType(*articleSpecificType).
			CreatedAt(mainPropositionData.CreatedAt).
			UpdatedAt(mainPropositionData.UpdatedAt).
			Build()
		if err != nil {
			log.Errorf("Error validating data for article %s of main proposition of voting %s: %s",
				mainPropositionData.Id, articleId, err.Error())
			return nil, err
		}

		mainProposition, err = proposition.NewBuilder().
			Id(mainPropositionData.Proposition.Id).
			Article(*articleDomain).
			Build()
		if err != nil {
			log.Errorf("Error validating data for main proposition %s (Article ID: %s) of voting %s: %s",
				mainPropositionData.Proposition.Id, mainPropositionData.Id, articleId, err.Error())
			return nil, err
		}
	}

	var relatedPropositions []proposition.Proposition
	for _, articleData := range articlesOfThePropositionsRelatedToVoting {
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
			log.Errorf("Error validating data for article %s related to voting %s: %s", articleData.Id,
				articleId, err.Error())
			return nil, err
		}

		relatedProposition, err := proposition.NewBuilder().
			Id(articleData.Proposition.Id).
			Article(*articleDomain).
			Build()
		if err != nil {
			log.Errorf("Error validating data for proposition %s of article %s related to voting %s: %s",
				articleData.Proposition.Id, articleData.Id, articleId, err)
			return nil, err
		}

		relatedPropositions = append(relatedPropositions, *relatedProposition)
	}

	var affectedPropositions []proposition.Proposition
	for _, articleData := range articlesOfThePropositionsAffectedByVoting {
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
			log.Errorf("Error validating data for article %s related to voting %s: %s", articleData.Id,
				votingArticle.Voting.Id, err.Error())
			return nil, err
		}

		affectedProposition, err := proposition.NewBuilder().
			Id(articleData.Proposition.Id).
			Article(*articleDomain).
			Build()
		if err != nil {
			log.Errorf("Error validating data for proposition %s of article %s affected by voting %s: %s",
				articleData.Proposition.Id, articleData.Id, articleId, err)
			return nil, err
		}

		affectedPropositions = append(affectedPropositions, *affectedProposition)
	}

	var articles []article.Article
	for _, articleData := range relatedArticles {
		articleType, err := articletype.NewBuilder().
			Id(articleData.ArticleType.Id).
			Description(articleData.ArticleType.Description).
			Codes(articleData.ArticleType.Codes).
			Color(articleData.ArticleType.Color).
			Build()
		if err != nil {
			log.Errorf("Error validating data for article type %s of article %s: %s",
				articleData.ArticleType.Id, articleData.Id, err.Error())
			return nil, err
		}

		articleBuilder := article.NewBuilder().
			Id(articleData.Id).
			AverageRating(articleData.AverageRating).
			NumberOfRatings(articleData.NumberOfRatings).
			Type(*articleType).
			CreatedAt(articleData.CreatedAt).
			UpdatedAt(articleData.UpdatedAt)

		if _, exists := userArticles[articleData.Id]; exists {
			articleBuilder.UserRating(userArticles[articleData.Id].Rating).
				ViewLater(userArticles[articleData.Id].ViewLater)
		}

		var articleDomain *article.Article
		var articleErr error
		if articleData.Event.Id != uuid.Nil {
			if articleData.Event.VideoUrl != "" {
				articleBuilder.MultimediaUrl(articleData.Event.VideoUrl)
			}

			articleSituationBuilder := articlesituation.NewBuilder()

			if !articleData.Event.EndsAt.IsZero() {
				articleSituationBuilder.EndsAt(articleData.Event.EndsAt)
			}

			articleSituation, err := articleSituationBuilder.
				Id(articleData.Event.EventSituation.Id).
				Description(articleData.Event.EventSituation.Description).
				Color(articleData.Event.EventSituation.Color).
				StartsAt(articleData.Event.StartsAt).
				Build()
			if err != nil {
				log.Errorf("Error validating data for article/event situation %s of event %s of article %s: %s",
					articleData.Event.EventSituation.Id, articleData.Event.Id, articleData.Id, err.Error())
				return nil, err
			}

			articleSpecificType, err := articletype.NewBuilder().
				Id(articleData.Event.EventType.Id).
				Description(articleData.Event.EventType.Description).
				Color(articleData.Event.EventType.Color).
				Build()
			if err != nil {
				log.Errorf("Error validating data for event type %s of event %s of article %s: %s",
					articleData.Event.EventType.Id, articleData.Event.Id, articleData.Id, err.Error())
				return nil, err
			}

			articleDomain, articleErr = articleBuilder.
				Title(articleData.Event.Title).
				Content(articleData.Event.Description).
				Situation(*articleSituation).
				SpecificType(*articleSpecificType).
				Build()
		} else {
			articleDomain, articleErr = articleBuilder.
				Title(fmt.Sprint("Boletim do dia ",
					articleData.Newsletter.ReferenceDate.Format("02/01/2006"))).
				Content(articleData.Newsletter.Description).
				Build()
		}
		if articleErr != nil {
			log.Errorf("Error validating data for article %s: %s", articleData.Id, articleErr.Error())
			return nil, articleErr
		}

		articles = append(articles, *articleDomain)
	}

	articleBuilder := article.NewBuilder()

	if _, exists := userArticles[articleId]; exists {
		articleBuilder.UserRating(userArticles[articleId].Rating).ViewLater(userArticles[articleId].ViewLater)
	}

	articleType, err := articletype.NewBuilder().
		Id(votingArticle.ArticleType.Id).
		Description(votingArticle.ArticleType.Description).
		Codes(votingArticle.ArticleType.Codes).
		Color(votingArticle.ArticleType.Color).
		Build()
	if err != nil {
		log.Errorf("Error validating data for article type %s of article %s: %s",
			votingArticle.ArticleType.Id, articleId, err.Error())
		return nil, err
	}

	articleDomain, err := articleBuilder.
		Id(votingArticle.Id).
		AverageRating(votingArticle.AverageRating).
		NumberOfRatings(votingArticle.NumberOfRatings).
		Type(*articleType).
		CreatedAt(votingArticle.CreatedAt).
		UpdatedAt(votingArticle.UpdatedAt).
		Build()
	if err != nil {
		log.Errorf("Error validating data for article %s of voting %s: %s", votingArticle.Id, articleId,
			err.Error())
		return nil, err
	}

	votingBuilder := voting.NewBuilder()

	if mainProposition != nil {
		votingBuilder.MainProposition(*mainProposition)
	}

	votingDomain, err := votingBuilder.Id(votingArticle.Voting.Id).
		Title(fmt.Sprint("Votação ", votingArticle.Voting.Code)).
		Description(votingArticle.Voting.Description).
		Result(votingArticle.Voting.Result).
		ResultAnnouncedAt(votingArticle.Voting.ResultAnnouncedAt).
		IsApproved(votingArticle.Voting.IsApproved).
		LegislativeBody(*legislativeBody).
		RelatedPropositions(relatedPropositions).
		AffectedPropositions(affectedPropositions).
		Article(*articleDomain).
		RelatedArticles(articles).
		Build()
	if err != nil {
		log.Errorf("Error validating data for voting %s of article %s: %s", votingArticle.Voting.Id, articleId,
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

	return votingDomain, nil
}
