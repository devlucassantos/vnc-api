package handlers

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"strings"
	"time"
	"vnc-api/api/endpoints/dto/request"
	"vnc-api/api/endpoints/dto/response"
	"vnc-api/api/endpoints/handlers/utils"
	"vnc-api/core/filters"
	"vnc-api/core/interfaces/services"
)

type Article struct {
	articleService     services.Article
	resourceService    services.Resources
	propositionService services.Proposition
	votingService      services.Voting
	eventService       services.Event
	newsletterService  services.Newsletter
}

func NewArticleHandler(articleService services.Article, resourceService services.Resources,
	propositionService services.Proposition, votingService services.Voting, eventService services.Event,
	newsletterService services.Newsletter) *Article {
	return &Article{
		articleService:     articleService,
		resourceService:    resourceService,
		propositionService: propositionService,
		votingService:      votingService,
		eventService:       eventService,
		newsletterService:  newsletterService,
	}
}

// GetArticles
// @ID          GetArticles
// @Summary     List most recent articles
// @Tags        Articles
// @Description This request is responsible for listing the most recent articles available on the platform.
// @Security    BearerAuth
// @Produce     json
// @Param       typeId                      query string false "Article type ID"
// @Param       specificTypeId              query string false "Article specific type ID"
// @Param       content                     query string false "Part of the content of the articles, in the title or content"
// @Param       startDate                   query string false "Date from which the articles were created. Accepted format: YYYY-MM-DD"
// @Param       endDate                     query string false "Date until which the articles were created. Accepted format: YYYY-MM-DD"
// @Param       propositionDeputyId         query string false "ID of the deputy who drafted the proposition"
// @Param       propositionPartyId          query string false "ID of the party that drafted the proposition"
// @Param       propositionExternalAuthorId query string false "ID of the external author who drafted the proposition"
// @Param       votingStartDate             query string false "Date from which the voting results were announced. Accepted format: YYYY-MM-DD"
// @Param       votingEndDate               query string false "Date until which the voting results were announced. Accepted format: YYYY-MM-DD"
// @Param       isVotingApproved            query bool   false "Is the voting approved?"
// @Param       votingLegislativeBodyId     query string false "ID of the legislative body responsible for the voting"
// @Param       eventStartDate              query string false "Date from which the events occurred. Accepted format: YYYY-MM-DD"
// @Param       eventEndDate                query string false "Date until which the events occurred. Accepted format: YYYY-MM-DD"
// @Param       eventSituationId            query string false "ID of the event situation"
// @Param       eventLegislativeBodyId      query string false "ID of the legislative body responsible for the event"
// @Param       eventRapporteurId           query string false "ID of the rapporteur (deputy) for one or more items on the event agenda"
// @Param       removeEventsInTheFuture     query bool   false "Remove events in the future?"
// @Param       page                        query int    false "Page number. By default, it is 1"
// @Param       itemsPerPage                query int    false "Number of articles returned per page. The default is 15 and the allowed values are between 1 and 100"
// @Success 200 {object} swagger.ArticlePagination "Successful request"
// @Failure 400 {object} swagger.HttpError         "Badly formatted request"
// @Failure 401 {object} swagger.HttpError         "Unauthorized access"
// @Failure 422 {object} swagger.HttpError         "Some of the data provided is invalid"
// @Failure 500 {object} swagger.HttpError         "An unexpected error occurred while processing the request"
// @Failure 503 {object} swagger.HttpError         "Some of the services/resources are temporarily unavailable"
// @Router /articles [GET]
func (instance Article) GetArticles(context echo.Context) error {
	userId := utils.GetUserIdFromAuthorizationHeader(context)

	articleFilter, httpError := getArticleQueryParametersFromContext(context)
	if httpError != nil {
		log.Warn("getArticleQueryParametersFromContext(): ", httpError.Message)
		return context.JSON(httpError.Code, httpError)
	}

	articleSlice, totalNumberOfArticles, err := instance.articleService.GetArticles(*articleFilter, userId)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			log.Error("Database unavailable: ", err.Error())
			return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
		}

		log.Error("Error retrieving articles: ", err.Error())
		return context.JSON(http.StatusInternalServerError, response.NewInternalServerError())
	}

	articles := make([]response.Article, 0)
	for _, articleData := range articleSlice {
		articles = append(articles, *response.NewArticle(articleData))
	}

	requestResult := response.Pagination{
		Page:         articleFilter.Pagination.GetPage(),
		ItemsPerPage: articleFilter.Pagination.GetItemsPerPage(),
		Total:        totalNumberOfArticles,
		Data:         articles,
	}

	return context.JSON(http.StatusOK, requestResult)
}

func getArticleQueryParametersFromContext(context echo.Context) (*filters.Article, *response.HttpError) {
	var articleFilter filters.Article

	typeIdParameter := context.QueryParam("typeId")
	if typeIdParameter != "" {
		parameter, parameterDescription := "typeId", "Article type ID"
		typeId, httpError := utils.ConvertFromStringToUuid(typeIdParameter, parameter, parameterDescription)
		if httpError != nil {
			log.Warn("Error converting the typeId parameter: ", httpError.Message)
			return nil, httpError
		}
		articleFilter.TypeId = &typeId
	}

	specificTypeIdParameter := context.QueryParam("specificTypeId")
	if specificTypeIdParameter != "" {
		parameter, parameterDescription := "specificTypeId", "Article specific type ID"
		specificTypeId, httpError := utils.ConvertFromStringToUuid(specificTypeIdParameter, parameter,
			parameterDescription)
		if httpError != nil {
			log.Warn("Error converting the specificTypeId parameter: ", httpError.Message)
			return nil, httpError
		}
		articleFilter.SpecificTypeId = &specificTypeId
	}

	articleFilter.Content = context.QueryParam("content")

	startDateParameter := context.QueryParam("startDate")
	if startDateParameter != "" {
		parameter, parameterDescription := "startDate", "Article start date"
		startDate, httpError := utils.ConvertFromStringToTime(startDateParameter, parameter, parameterDescription)
		if httpError != nil {
			log.Warn("Error converting the startDate parameter: ", httpError.Message)
			return nil, httpError
		}
		articleFilter.StartDate = &startDate
	}

	endDateParameter := context.QueryParam("endDate")
	if endDateParameter != "" {
		parameter, parameterDescription := "endDate", "Article end date"
		endDate, httpError := utils.ConvertFromStringToTime(endDateParameter, parameter, parameterDescription)
		if httpError != nil {
			log.Warn("Error converting the endDate parameter: ", httpError.Message)
			return nil, httpError
		}
		articleFilter.EndDate = &endDate
	}

	if articleFilter.StartDate != nil && articleFilter.EndDate != nil &&
		articleFilter.StartDate.After(*articleFilter.EndDate) {
		errorMessage := fmt.Sprint("Invalid parameters: The article start date parameter (startDate) cannot be " +
			"greater than the article end date parameter (endDate)")
		log.Warn(errorMessage)
		return nil, response.NewHttpError(http.StatusUnprocessableEntity, errorMessage)
	}

	propositionDeputyIdParameter := context.QueryParam("propositionDeputyId")
	if propositionDeputyIdParameter != "" {
		parameter, parameterDescription := "propositionDeputyId", "Proposition deputy ID"
		propositionDeputyId, httpError := utils.ConvertFromStringToUuid(propositionDeputyIdParameter, parameter,
			parameterDescription)
		if httpError != nil {
			log.Warn("Error converting the propositionDeputyId parameter: ", httpError.Message)
			return nil, httpError
		}
		articleFilter.Proposition.DeputyId = &propositionDeputyId
	}

	propositionPartyIdParameter := context.QueryParam("propositionPartyId")
	if propositionPartyIdParameter != "" {
		parameter, parameterDescription := "propositionPartyId", "Proposition party ID"
		propositionPartyId, httpError := utils.ConvertFromStringToUuid(propositionPartyIdParameter, parameter,
			parameterDescription)
		if httpError != nil {
			log.Warn("Error converting the propositionPartyId parameter: ", httpError.Message)
			return nil, httpError
		}
		articleFilter.Proposition.PartyId = &propositionPartyId
	}

	propositionExternalAuthorIdParameter := context.QueryParam("propositionExternalAuthorId")
	if propositionExternalAuthorIdParameter != "" {
		parameter, parameterDescription := "propositionExternalAuthorId", "Proposition external author ID"
		propositionExternalAuthorId, httpError := utils.ConvertFromStringToUuid(propositionExternalAuthorIdParameter,
			parameter, parameterDescription)
		if httpError != nil {
			log.Warn("Error converting the propositionExternalAuthorId parameter: ", httpError.Message)
			return nil, httpError
		}
		articleFilter.Proposition.ExternalAuthorId = &propositionExternalAuthorId
	}

	votingStartDateParameter := context.QueryParam("votingStartDate")
	if votingStartDateParameter != "" {
		parameter, parameterDescription := "votingStartDate", "Voting start date"
		votingStartDate, httpError := utils.ConvertFromStringToTime(votingStartDateParameter, parameter,
			parameterDescription)
		if httpError != nil {
			log.Warn("Error converting the votingStartDate parameter: ", httpError.Message)
			return nil, httpError
		}
		articleFilter.Voting.StartDate = &votingStartDate
	}

	votingEndDateParameter := context.QueryParam("votingEndDate")
	if votingEndDateParameter != "" {
		parameter, parameterDescription := "votingEndDate", "Voting end date"
		votingEndDate, httpError := utils.ConvertFromStringToTime(votingEndDateParameter, parameter,
			parameterDescription)
		if httpError != nil {
			log.Warn("Error converting the votingEndDate parameter: ", httpError.Message)
			return nil, httpError
		}
		articleFilter.Voting.EndDate = &votingEndDate
	}

	if articleFilter.Voting.StartDate != nil && articleFilter.Voting.EndDate != nil &&
		articleFilter.Voting.StartDate.After(*articleFilter.Voting.EndDate) {
		errorMessage := fmt.Sprint("Invalid parameters: The voting start date parameter (votingStartDate) " +
			"cannot be greater than the voting end date parameter (votingEndDate)")
		log.Warn(errorMessage)
		return nil, response.NewHttpError(http.StatusUnprocessableEntity, errorMessage)
	}

	isVotingApprovedParameter := context.QueryParam("isVotingApproved")
	if isVotingApprovedParameter != "" {
		parameter, parameterDescription := "isVotingApproved", "Is the voting approved?"
		isVotingApproved, httpError := utils.ConvertFromStringToBool(isVotingApprovedParameter, parameter,
			parameterDescription)
		if httpError != nil {
			log.Warn("Error converting the isVotingApproved parameter: ", httpError.Message)
			return nil, httpError
		}
		articleFilter.Voting.IsVotingApproved = &isVotingApproved
	}

	votingLegislativeBodyIdParameter := context.QueryParam("votingLegislativeBodyId")
	if votingLegislativeBodyIdParameter != "" {
		parameter, parameterDescription := "votingLegislativeBodyId", "Voting Legislative Body ID"
		votingLegislativeBodyId, httpError := utils.ConvertFromStringToUuid(votingLegislativeBodyIdParameter, parameter,
			parameterDescription)
		if httpError != nil {
			log.Warn("Error converting the votingLegislativeBodyId parameter: ", httpError.Message)
			return nil, httpError
		}
		articleFilter.Voting.LegislativeBodyId = &votingLegislativeBodyId
	}

	eventStartDateParameter := context.QueryParam("eventStartDate")
	if eventStartDateParameter != "" {
		parameter, parameterDescription := "eventStartDate", "Event start date"
		eventStartDate, httpError := utils.ConvertFromStringToTime(eventStartDateParameter, parameter,
			parameterDescription)
		if httpError != nil {
			log.Warn("Error converting the eventStartDate parameter: ", httpError.Message)
			return nil, httpError
		}
		articleFilter.Event.StartDate = &eventStartDate
	}

	eventEndDateParameter := context.QueryParam("eventEndDate")
	if eventEndDateParameter != "" {
		parameter, parameterDescription := "eventEndDate", "Event end date"
		eventEndDate, httpError := utils.ConvertFromStringToTime(eventEndDateParameter, parameter, parameterDescription)
		if httpError != nil {
			log.Warn("Error converting the eventEndDate parameter: ", httpError.Message)
			return nil, httpError
		}
		articleFilter.Event.EndDate = &eventEndDate
	}

	if articleFilter.Event.StartDate != nil && articleFilter.Event.EndDate != nil &&
		articleFilter.Event.StartDate.After(*articleFilter.Event.EndDate) {
		errorMessage := fmt.Sprint("Invalid parameters: The event start date parameter (eventStartDate) cannot " +
			"be greater than the event end date parameter (eventEndDate)")
		log.Warn(errorMessage)
		return nil, response.NewHttpError(http.StatusUnprocessableEntity, errorMessage)
	}

	eventSituationIdParameter := context.QueryParam("eventSituationId")
	if eventSituationIdParameter != "" {
		parameter, parameterDescription := "eventSituationId", "Event Situation ID"
		eventSituationId, httpError := utils.ConvertFromStringToUuid(eventSituationIdParameter, parameter,
			parameterDescription)
		if httpError != nil {
			log.Warn("Error converting the eventSituationId parameter: ", httpError.Message)
			return nil, httpError
		}
		articleFilter.Event.SituationId = &eventSituationId
	}

	eventLegislativeBodyIdParameter := context.QueryParam("eventLegislativeBodyId")
	if eventLegislativeBodyIdParameter != "" {
		parameter, parameterDescription := "eventLegislativeBodyId", "Event Legislative Body ID"
		eventLegislativeBodyId, httpError := utils.ConvertFromStringToUuid(eventLegislativeBodyIdParameter, parameter,
			parameterDescription)
		if httpError != nil {
			log.Warn("Error converting the eventLegislativeBodyId parameter: ", httpError.Message)
			return nil, httpError
		}
		articleFilter.Event.LegislativeBodyId = &eventLegislativeBodyId
	}

	eventRapporteurIdParameter := context.QueryParam("eventRapporteurId")
	if eventRapporteurIdParameter != "" {
		parameter, parameterDescription := "eventRapporteurId", "Event Rapporteur ID"
		eventRapporteurId, httpError := utils.ConvertFromStringToUuid(eventRapporteurIdParameter, parameter,
			parameterDescription)
		if httpError != nil {
			log.Warn("Error converting the eventRapporteurId parameter: ", httpError.Message)
			return nil, httpError
		}
		articleFilter.Event.RapporteurId = &eventRapporteurId
	}

	removeEventsInTheFutureParameter := context.QueryParam("removeEventsInTheFuture")
	if removeEventsInTheFutureParameter != "" {
		parameter, parameterDescription := "removeEventsInTheFuture", "Remove events in the future?"
		removeEventsInTheFuture, httpError := utils.ConvertFromStringToBool(removeEventsInTheFutureParameter, parameter,
			parameterDescription)
		if httpError != nil {
			log.Warn("Error converting the removeEventsInTheFuture parameter: ", httpError.Message)
			return nil, httpError
		}
		articleFilter.Event.RemoveEventsInTheFuture = &removeEventsInTheFuture

		if removeEventsInTheFuture {
			eventEndDate := time.Now()
			articleFilter.Event.EndDate = &eventEndDate
		}
	}

	err := articleFilter.HasConflict()
	if err != nil {
		log.Warn(err.Error())
		return nil, response.NewHttpError(http.StatusUnprocessableEntity, err.Error())
	}

	pageParameter := context.QueryParam("page")
	if pageParameter != "" {
		parameter, parameterDescription := "page", "Page"
		page, httpError := utils.ConvertFromStringToInt(pageParameter, parameter, parameterDescription)
		if httpError != nil {
			log.Warn("Error converting the page parameter: ", httpError.Message)
			return nil, httpError
		}
		articleFilter.Pagination.Page = &page
	}

	itemsPerPageParameter := context.QueryParam("itemsPerPage")
	if itemsPerPageParameter != "" {
		parameter, parameterDescription := "itemsPerPage", "Items per page"
		itemsPerPage, httpError := utils.ConvertFromStringToInt(itemsPerPageParameter, parameter, parameterDescription)
		if httpError != nil {
			log.Warn("Error converting the itemsPerPage parameter: ", httpError.Message)
			return nil, httpError
		}

		if itemsPerPage > 100 {
			errorMessage := fmt.Sprint("Invalid parameter: Items per page (itemsPerPage) must be less than or " +
				"equal to 100")
			log.Warnf("Parameter out of allowed range: %s (Value: %d)", errorMessage, itemsPerPage)
			return nil, response.NewHttpError(http.StatusUnprocessableEntity, errorMessage)
		}

		articleFilter.Pagination.ItemsPerPage = &itemsPerPage
	}

	return &articleFilter, nil
}

// GetTrendingArticles
// @ID          GetTrendingArticles
// @Summary     List trending articles
// @Tags        Articles
// @Description This request is responsible for listing trending articles on the platform.
// @Security    BearerAuth
// @Produce     json
// @Param       typeId                      query string false "Article type ID"
// @Param       specificTypeId              query string false "Article specific type ID"
// @Param       content                     query string false "Part of the content of the articles, in the title or content"
// @Param       startDate                   query string false "Date from which the articles were created. Accepted format: YYYY-MM-DD"
// @Param       endDate                     query string false "Date until which the articles were created. Accepted format: YYYY-MM-DD"
// @Param       propositionDeputyId         query string false "ID of the deputy who drafted the proposition"
// @Param       propositionPartyId          query string false "ID of the party that drafted the proposition"
// @Param       propositionExternalAuthorId query string false "ID of the external author who drafted the proposition"
// @Param       votingStartDate             query string false "Date from which the voting results were announced. Accepted format: YYYY-MM-DD"
// @Param       votingEndDate               query string false "Date until which the voting results were announced. Accepted format: YYYY-MM-DD"
// @Param       isVotingApproved            query bool   false "Is the voting approved?"
// @Param       votingLegislativeBodyId     query string false "ID of the legislative body responsible for the voting"
// @Param       eventStartDate              query string false "Date from which the events occurred. Accepted format: YYYY-MM-DD"
// @Param       eventEndDate                query string false "Date until which the events occurred. Accepted format: YYYY-MM-DD"
// @Param       eventSituationId            query string false "ID of the event situation"
// @Param       eventLegislativeBodyId      query string false "ID of the legislative body responsible for the event"
// @Param       eventRapporteurId           query string false "ID of the rapporteur (deputy) for one or more items on the event agenda"
// @Param       removeEventsInTheFuture     query bool   false "Remove events in the future?"
// @Param       page                        query int    false "Page number. By default, it is 1"
// @Param       itemsPerPage                query int    false "Number of articles returned per page. The default is 15 and the allowed values are between 1 and 100"
// @Success 200 {object} swagger.ArticlePagination "Successful request"
// @Failure 400 {object} swagger.HttpError         "Badly formatted request"
// @Failure 401 {object} swagger.HttpError         "Unauthorized access"
// @Failure 422 {object} swagger.HttpError         "Some of the data provided is invalid"
// @Failure 500 {object} swagger.HttpError         "An unexpected error occurred while processing the request"
// @Failure 503 {object} swagger.HttpError         "Some of the services/resources are temporarily unavailable"
// @Router /articles/trending [GET]
func (instance Article) GetTrendingArticles(context echo.Context) error {
	userId := utils.GetUserIdFromAuthorizationHeader(context)

	articleFilter, httpError := getArticleQueryParametersFromContext(context)
	if httpError != nil {
		log.Warn("getArticleQueryParametersFromContext(): ", httpError.Message)
		return context.JSON(httpError.Code, httpError)
	}

	articleSlice, totalNumberOfArticles, err := instance.articleService.GetTrendingArticles(*articleFilter, userId)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			log.Error("Database unavailable: ", err.Error())
			return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
		}

		log.Error("Error retrieving trending articles: ", err.Error())
		return context.JSON(http.StatusInternalServerError, response.NewInternalServerError())
	}

	articles := make([]response.Article, 0)
	for _, articleData := range articleSlice {
		articles = append(articles, *response.NewArticle(articleData))
	}

	requestResult := response.Pagination{
		Page:         articleFilter.Pagination.GetPage(),
		ItemsPerPage: articleFilter.Pagination.GetItemsPerPage(),
		Total:        totalNumberOfArticles,
		Data:         articles,
	}

	return context.JSON(http.StatusOK, requestResult)
}

// GetTrendingArticlesByType
// @ID          GetTrendingArticlesByType
// @Summary     List trending articles by article types
// @Tags        Articles
// @Description This request is responsible for listing the trending articles by article types or specific types.
// @Security    BearerAuth
// @Produce     json
// @Param       articleTypeIds         query string false "List of IDs of the types of articles that should be returned (separated by commas). By default, it returns all types"
// @Param       articleSpecificTypeIds query string false "List of IDs of the specific types of articles that should be returned (separated by commas). By default, it returns all specific types"
// @Param       itemsPerType           query int    false "Number of articles returned by type. The default is 5 and the allowed values are between 1 and 20"
// @Success 200 {object} swagger.ArticleTypeWithSpecificTypes "Successful request"
// @Failure 400 {object} swagger.HttpError                    "Badly formatted request"
// @Failure 401 {object} swagger.HttpError                    "Unauthorized access"
// @Failure 500 {object} swagger.HttpError                    "An unexpected error occurred while processing the request"
// @Failure 503 {object} swagger.HttpError                    "Some of the services/resources are temporarily unavailable"
// @Router /articles/trending/type [GET]
func (instance Article) GetTrendingArticlesByType(context echo.Context) error {
	userId := utils.GetUserIdFromAuthorizationHeader(context)

	var articleTypeIds []uuid.UUID
	articleTypeIdsParameter := context.QueryParam("articleTypeIds")
	if articleTypeIdsParameter != "" {
		articleTypeIdsAsStringSlice := strings.Split(articleTypeIdsParameter, ",")
		for index, articleTypeIdAsString := range articleTypeIdsAsStringSlice {
			parameter, parameterDescription := "articleTypeIds", fmt.Sprintf("%d°th article type ID", index+1)
			articleTypeId, httpError := utils.ConvertFromStringToUuid(articleTypeIdAsString, parameter,
				parameterDescription)
			if httpError != nil {
				log.Warn("Error converting the articleTypeIds parameter: ", httpError.Message)
				return context.JSON(httpError.Code, httpError)
			}
			articleTypeIds = append(articleTypeIds, articleTypeId)
		}
	}

	var articleSpecificTypeIds []uuid.UUID
	articleSpecificTypeIdsParameter := context.QueryParam("articleSpecificTypeIds")
	if articleSpecificTypeIdsParameter != "" {
		articleSpecificTypeIdsAsStringSlice := strings.Split(articleSpecificTypeIdsParameter, ",")
		for index, articleSpecificTypeIdAsString := range articleSpecificTypeIdsAsStringSlice {
			parameter := "articleSpecificTypeIds"
			parameterDescription := fmt.Sprintf("%d°th article specific type ID", index+1)
			articleSpecificTypeId, httpError := utils.ConvertFromStringToUuid(articleSpecificTypeIdAsString, parameter,
				parameterDescription)
			if httpError != nil {
				log.Warn("Error converting the articleSpecificTypeIds parameter: ", httpError.Message)
				return context.JSON(httpError.Code, httpError)
			}
			articleSpecificTypeIds = append(articleSpecificTypeIds, articleSpecificTypeId)
		}
	}

	itemsPerType := 5
	itemsPerTypeParameter := context.QueryParam("itemsPerType")
	if itemsPerTypeParameter != "" {
		var httpError *response.HttpError
		parameter, parameterDescription := "itemsPerType", "Items per type"
		itemsPerType, httpError = utils.ConvertFromStringToInt(itemsPerTypeParameter, parameter, parameterDescription)
		if httpError != nil {
			log.Warn("Error converting the itemsPerType parameter: ", httpError.Message)
			return context.JSON(httpError.Code, httpError)
		}

		if itemsPerType > 20 {
			errorMessage := fmt.Sprint("Invalid parameter: Items per type (itemsPerType) must be less than or " +
				"equal to 20")
			log.Warnf("Parameter out of allowed range: %s (Value: %d)", errorMessage, itemsPerType)
			return context.JSON(http.StatusUnprocessableEntity, response.NewHttpError(http.StatusUnprocessableEntity,
				errorMessage))
		}
	}

	articleTypeSlice := make([]response.ArticleType, 0)
	if articleSpecificTypeIds != nil {
		for _, articleSpecificTypeId := range articleSpecificTypeIds {
			trendingArticles, err := instance.articleService.GetTrendingArticlesBySpecificTypeId(articleSpecificTypeId,
				itemsPerType, userId)
			if err != nil {
				if strings.Contains(err.Error(), "connection refused") {
					log.Error("Database unavailable: ", err.Error())
					return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
				}

				log.Errorf("Error retrieving trending articles by specific type %s: %s", articleSpecificTypeId,
					err.Error())
				return context.JSON(http.StatusInternalServerError, response.NewInternalServerError())
			}

			articleTypeSlice = response.SortingArticleTypesWithSpecificTypesAndArticles(trendingArticles,
				articleTypeSlice)
		}

		if articleTypeIds != nil {
			var articleTypesWithoutArticles []uuid.UUID
			for _, articleTypeId := range articleTypeIds {
				var articleTypeHasArticles bool
				for _, articleType := range articleTypeSlice {
					if articleType.Id == articleTypeId {
						articleTypeHasArticles = true
						break
					}
				}

				if !articleTypeHasArticles {
					articleTypesWithoutArticles = append(articleTypesWithoutArticles, articleTypeId)
				}
			}

			for _, articleTypeId := range articleTypesWithoutArticles {
				trendingArticles, err := instance.articleService.GetTrendingArticlesByTypeId(articleTypeId,
					itemsPerType, userId)
				if err != nil {
					if strings.Contains(err.Error(), "connection refused") {
						log.Error("Database unavailable: ", err.Error())
						return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
					}

					log.Errorf("Error retrieving trending articles by type %s: %s", articleTypeId, err.Error())
					return context.JSON(http.StatusInternalServerError, response.NewInternalServerError())
				}

				articleTypeSlice = response.SortingArticleTypeWithArticles(trendingArticles, articleTypeSlice)
			}
		}
	} else if articleTypeIds != nil {
		for _, articleTypeId := range articleTypeIds {
			trendingArticles, err := instance.articleService.GetTrendingArticlesByTypeId(articleTypeId, itemsPerType,
				userId)
			if err != nil {
				if strings.Contains(err.Error(), "connection refused") {
					log.Error("Database unavailable: ", err.Error())
					return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
				}

				log.Errorf("Error retrieving trending articles by type %s: %s", articleTypeId, err.Error())
				return context.JSON(http.StatusInternalServerError, response.NewInternalServerError())
			}

			articleTypeSlice = response.SortingArticleTypeWithArticles(trendingArticles, articleTypeSlice)
		}
	} else {
		articleTypes, err := instance.resourceService.GetArticleTypes()
		if err != nil {
			if strings.Contains(err.Error(), "connection refused") {
				log.Error("Database unavailable: ", err.Error())
				return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
			}

			log.Errorf("Error retrieving article types data from the database: %s", err.Error())
			return context.JSON(http.StatusInternalServerError, response.NewInternalServerError())
		}

		for _, articleType := range articleTypes {
			if strings.Contains(articleType.Codes(), "proposition") {
				propositionTypes, err := instance.resourceService.GetPropositionTypes()
				if err != nil {
					if strings.Contains(err.Error(), "connection refused") {
						log.Error("Database unavailable: ", err.Error())
						return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
					}

					log.Errorf("Error retrieving proposition types data from the database: %s", err.Error())
					return context.JSON(http.StatusInternalServerError, response.NewInternalServerError())
				}

				for _, propositionType := range propositionTypes {
					trendingArticles, err := instance.articleService.GetTrendingArticlesBySpecificTypeId(
						propositionType.Id(), itemsPerType, userId)
					if err != nil {
						if strings.Contains(err.Error(), "connection refused") {
							log.Error("Database unavailable: ", err.Error())
							return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
						}

						log.Errorf("Error retrieving trending articles by specific type %s: %s",
							propositionType.Id(), err.Error())
						return context.JSON(http.StatusInternalServerError, response.NewInternalServerError())
					}

					articleTypeSlice = response.SortingArticleTypesWithSpecificTypesAndArticles(trendingArticles,
						articleTypeSlice)
				}
			} else if strings.Contains(articleType.Codes(), "event") {
				eventTypes, err := instance.resourceService.GetEventTypes()
				if err != nil {
					if strings.Contains(err.Error(), "connection refused") {
						log.Error("Database unavailable: ", err.Error())
						return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
					}

					log.Errorf("Error retrieving event types data from the database: %s", err.Error())
					return context.JSON(http.StatusInternalServerError, response.NewInternalServerError())
				}

				for _, eventType := range eventTypes {
					trendingArticles, err := instance.articleService.GetTrendingArticlesBySpecificTypeId(eventType.Id(),
						itemsPerType, userId)
					if err != nil {
						if strings.Contains(err.Error(), "connection refused") {
							log.Error("Database unavailable: ", err.Error())
							return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
						}

						log.Errorf("Error retrieving trending articles by specific type %s: %s", eventType.Id(),
							err.Error())
						return context.JSON(http.StatusInternalServerError, response.NewInternalServerError())
					}

					articleTypeSlice = response.SortingArticleTypesWithSpecificTypesAndArticles(trendingArticles,
						articleTypeSlice)
				}
			} else {
				trendingArticles, err := instance.articleService.GetTrendingArticlesByTypeId(articleType.Id(),
					itemsPerType, userId)
				if err != nil {
					if strings.Contains(err.Error(), "connection refused") {
						log.Error("Database unavailable: ", err.Error())
						return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
					}

					log.Errorf("Error retrieving trending articles by type %s: %s", articleType.Id(),
						err.Error())
					return context.JSON(http.StatusInternalServerError, response.NewInternalServerError())
				}

				articleTypeSlice = response.SortingArticleTypeWithArticles(trendingArticles, articleTypeSlice)
			}
		}
	}

	return context.JSON(http.StatusOK, articleTypeSlice)
}

// GetArticlesToViewLater
// @ID          GetArticlesToViewLater
// @Summary     List articles bookmarked for later viewing by the user
// @Tags        Articles
// @Description This request is responsible for listing the articles bookmarked for later viewing by the user on the platform. The articles will be listed in the order in which the user bookmarked the articles.
// @Security    BearerAuth
// @Produce     json
// @Param       typeId                      query string false "Article type ID"
// @Param       specificTypeId              query string false "Article specific type ID"
// @Param       content                     query string false "Part of the content of the articles, in the title or content"
// @Param       startDate                   query string false "Date from which the articles were created. Accepted format: YYYY-MM-DD"
// @Param       endDate                     query string false "Date until which the articles were created. Accepted format: YYYY-MM-DD"
// @Param       propositionDeputyId         query string false "ID of the deputy who drafted the proposition"
// @Param       propositionPartyId          query string false "ID of the party that drafted the proposition"
// @Param       propositionExternalAuthorId query string false "ID of the external author who drafted the proposition"
// @Param       votingStartDate             query string false "Date from which the voting results were announced. Accepted format: YYYY-MM-DD"
// @Param       votingEndDate               query string false "Date until which the voting results were announced. Accepted format: YYYY-MM-DD"
// @Param       isVotingApproved            query bool   false "Is the voting approved?"
// @Param       votingLegislativeBodyId     query string false "ID of the legislative body responsible for the voting"
// @Param       eventStartDate              query string false "Date from which the events occurred. Accepted format: YYYY-MM-DD"
// @Param       eventEndDate                query string false "Date until which the events occurred. Accepted format: YYYY-MM-DD"
// @Param       eventSituationId            query string false "ID of the event situation"
// @Param       eventLegislativeBodyId      query string false "ID of the legislative body responsible for the event"
// @Param       eventRapporteurId           query string false "ID of the rapporteur (deputy) for one or more items on the event agenda"
// @Param       removeEventsInTheFuture     query bool   false "Remove events in the future?"
// @Param       page                        query int    false "Page number. By default, it is 1"
// @Param       itemsPerPage                query int    false "Number of articles returned per page. The default is 15 and the allowed values are between 1 and 100"
// @Success 200 {object} swagger.ArticlePagination "Successful request"
// @Failure 400 {object} swagger.HttpError         "Badly formatted request"
// @Failure 401 {object} swagger.HttpError         "Unauthorized access"
// @Failure 422 {object} swagger.HttpError         "Some of the data provided is invalid"
// @Failure 500 {object} swagger.HttpError         "An unexpected error occurred while processing the request"
// @Failure 503 {object} swagger.HttpError         "Some of the services/resources are temporarily unavailable"
// @Router /articles/view-later [GET]
func (instance Article) GetArticlesToViewLater(context echo.Context) error {
	userId := utils.GetUserIdFromAuthorizationHeader(context)

	articleFilter, httpError := getArticleQueryParametersFromContext(context)
	if httpError != nil {
		log.Warn("getArticleQueryParametersFromContext(): ", httpError.Message)
		return context.JSON(httpError.Code, httpError)
	}

	articleSlice, totalNumberOfArticles, err := instance.articleService.GetArticlesToViewLater(*articleFilter, userId)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			log.Error("Database unavailable: ", err.Error())
			return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
		}

		log.Errorf("Error retrieving articles bookmarked for later viewing by user %s: %s", userId, err.Error())
		return context.JSON(http.StatusInternalServerError, response.NewInternalServerError())
	}

	articles := make([]response.Article, 0)
	for _, articleData := range articleSlice {
		articles = append(articles, *response.NewArticle(articleData))
	}

	requestResult := response.Pagination{
		Page:         articleFilter.Pagination.GetPage(),
		ItemsPerPage: articleFilter.Pagination.GetItemsPerPage(),
		Total:        totalNumberOfArticles,
		Data:         articles,
	}

	return context.JSON(http.StatusOK, requestResult)
}

// GetPropositionArticleById
// @ID          GetPropositionArticleById
// @Summary     Get article details by ID (Only for proposition articles)
// @Tags        Articles
// @Description This request is responsible for looking up the details of an article of a proposition by the article ID.
// @Security    BearerAuth
// @Produce     json
// @Param       articleId path string true "Article ID"
// @Success 200 {object} swagger.PropositionArticle "Successful request"
// @Failure 400 {object} swagger.HttpError          "Badly formatted request"
// @Failure 401 {object} swagger.HttpError          "Unauthorized access"
// @Failure 404 {object} swagger.HttpError          "Requested resource not found"
// @Failure 422 {object} swagger.HttpError          "Some of the data provided is invalid"
// @Failure 500 {object} swagger.HttpError          "An unexpected error occurred while processing the request"
// @Failure 503 {object} swagger.HttpError          "Some of the services/resources are temporarily unavailable"
// @Router /articles/{articleId}/proposition [GET]
func (instance Article) GetPropositionArticleById(context echo.Context) error {
	userId := utils.GetUserIdFromAuthorizationHeader(context)

	articleIdParameter := context.Param("articleId")
	parameter, parameterDescription := "articleId", "Article ID"
	articleId, httpError := utils.ConvertFromStringToUuid(articleIdParameter, parameter, parameterDescription)
	if httpError != nil {
		log.Warn("Error converting the articleId parameter: ", httpError.Message)
		return context.JSON(httpError.Code, httpError)
	}

	propositionData, err := instance.propositionService.GetPropositionByArticleId(articleId, userId)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			log.Warnf("Proposition article %s could not be found: %s", articleId, err.Error())
			return context.JSON(http.StatusNotFound, response.NewHttpError(http.StatusNotFound,
				"Proposition article not found"))
		} else if strings.Contains(err.Error(), "connection refused") {
			log.Error("Database unavailable: ", err.Error())
			return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
		}

		log.Errorf("Error retrieving proposition article %s: %s", articleId, err.Error())
		return context.JSON(http.StatusInternalServerError, response.NewInternalServerError())
	}

	propositionArticle := response.NewPropositionArticle(*propositionData)

	return context.JSON(http.StatusOK, propositionArticle)
}

// GetVotingArticleById
// @ID          GetVotingArticleById
// @Summary     Get article details by ID (Only for voting articles)
// @Tags        Articles
// @Description This request is responsible for looking up the details of an article of a voting by the article ID.
// @Security    BearerAuth
// @Produce     json
// @Param       articleId path string true "Article ID"
// @Success 200 {object} swagger.VotingArticle "Successful request"
// @Failure 400 {object} swagger.HttpError     "Badly formatted request"
// @Failure 401 {object} swagger.HttpError     "Unauthorized access"
// @Failure 404 {object} swagger.HttpError     "Requested resource not found"
// @Failure 422 {object} swagger.HttpError     "Some of the data provided is invalid"
// @Failure 500 {object} swagger.HttpError     "An unexpected error occurred while processing the request"
// @Failure 503 {object} swagger.HttpError     "Some of the services/resources are temporarily unavailable"
// @Router /articles/{articleId}/voting [GET]
func (instance Article) GetVotingArticleById(context echo.Context) error {
	userId := utils.GetUserIdFromAuthorizationHeader(context)

	articleIdParameter := context.Param("articleId")
	parameter, parameterDescription := "articleId", "Article ID"
	articleId, httpError := utils.ConvertFromStringToUuid(articleIdParameter, parameter, parameterDescription)
	if httpError != nil {
		log.Warn("Error converting the articleId parameter: ", httpError.Message)
		return context.JSON(httpError.Code, httpError)
	}

	votingData, err := instance.votingService.GetVotingByArticleId(articleId, userId)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			log.Warnf("Voting article %s could not be found: %s", articleId, err.Error())
			return context.JSON(http.StatusNotFound, response.NewHttpError(http.StatusNotFound,
				"Voting article not found"))
		} else if strings.Contains(err.Error(), "connection refused") {
			log.Error("Database unavailable: ", err.Error())
			return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
		}

		log.Errorf("Error retrieving voting article %s: %s", articleId, err.Error())
		return context.JSON(http.StatusInternalServerError, response.NewInternalServerError())
	}

	votingArticle := response.NewVotingArticle(*votingData)

	return context.JSON(http.StatusOK, votingArticle)
}

// GetEventArticleById
// @ID          GetEventArticleById
// @Summary     Get article details by ID (Only for event articles)
// @Tags        Articles
// @Description This request is responsible for looking up the details of an article of an event by the article ID.
// @Security    BearerAuth
// @Produce     json
// @Param       articleId path string true "Article ID"
// @Success 200 {object} swagger.EventArticle "Successful request"
// @Failure 400 {object} swagger.HttpError    "Badly formatted request"
// @Failure 401 {object} swagger.HttpError    "Unauthorized access"
// @Failure 404 {object} swagger.HttpError    "Requested resource not found"
// @Failure 422 {object} swagger.HttpError    "Some of the data provided is invalid"
// @Failure 500 {object} swagger.HttpError    "An unexpected error occurred while processing the request"
// @Failure 503 {object} swagger.HttpError    "Some of the services/resources are temporarily unavailable"
// @Router /articles/{articleId}/event [GET]
func (instance Article) GetEventArticleById(context echo.Context) error {
	userId := utils.GetUserIdFromAuthorizationHeader(context)

	articleIdParameter := context.Param("articleId")
	parameter, parameterDescription := "articleId", "Article ID"
	articleId, httpError := utils.ConvertFromStringToUuid(articleIdParameter, parameter, parameterDescription)
	if httpError != nil {
		log.Warn("Error converting the articleId parameter: ", httpError.Message)
		return context.JSON(httpError.Code, httpError)
	}

	eventData, err := instance.eventService.GetEventByArticleId(articleId, userId)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			log.Warnf("Event article %s could not be found: %s", articleId, err.Error())
			return context.JSON(http.StatusNotFound, response.NewHttpError(http.StatusNotFound,
				"Event article not found"))
		} else if strings.Contains(err.Error(), "connection refused") {
			log.Error("Database unavailable: ", err.Error())
			return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
		}

		log.Errorf("Error retrieving event article %s: %s", articleId, err.Error())
		return context.JSON(http.StatusInternalServerError, response.NewInternalServerError())
	}

	eventArticle := response.NewEventArticle(*eventData)

	return context.JSON(http.StatusOK, eventArticle)
}

// GetNewsletterArticleById
// @ID          GetNewsletterArticleById
// @Summary     Get article details by ID (Only for newsletter articles)
// @Tags        Articles
// @Description This request is responsible for looking up the details of an article of a newsletter by the article ID.
// @Security    BearerAuth
// @Produce     json
// @Param       articleId path string true "Article ID"
// @Success 200 {object} swagger.NewsletterArticle "Successful request"
// @Failure 400 {object} swagger.HttpError         "Badly formatted request"
// @Failure 401 {object} swagger.HttpError         "Unauthorized access"
// @Failure 404 {object} swagger.HttpError         "Requested resource not found"
// @Failure 422 {object} swagger.HttpError         "Some of the data provided is invalid"
// @Failure 500 {object} swagger.HttpError         "An unexpected error occurred while processing the request"
// @Failure 503 {object} swagger.HttpError         "Some of the services/resources are temporarily unavailable"
// @Router /articles/{articleId}/newsletter [GET]
func (instance Article) GetNewsletterArticleById(context echo.Context) error {
	userId := utils.GetUserIdFromAuthorizationHeader(context)

	articleIdParameter := context.Param("articleId")
	parameter, parameterDescription := "articleId", "Article ID"
	articleId, httpError := utils.ConvertFromStringToUuid(articleIdParameter, parameter, parameterDescription)
	if httpError != nil {
		log.Warn("Error converting the articleId parameter: ", httpError.Message)
		return context.JSON(httpError.Code, httpError)
	}

	newsletterData, err := instance.newsletterService.GetNewsletterByArticleId(articleId, userId)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			log.Warnf("Newsletter article %s could not be found: %s", articleId, err.Error())
			return context.JSON(http.StatusNotFound, response.NewHttpError(http.StatusNotFound,
				"Newsletter article not found"))
		} else if strings.Contains(err.Error(), "connection refused") {
			log.Error("Database unavailable: ", err.Error())
			return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
		}

		log.Errorf("Error retrieving newsletter article %s: %s", articleId, err.Error())
		return context.JSON(http.StatusInternalServerError, response.NewInternalServerError())
	}

	newsletterArticle := response.NewNewsletterArticle(*newsletterData)

	return context.JSON(http.StatusOK, newsletterArticle)
}

// SaveArticleRating
// @ID          SaveArticleRating
// @Summary     Rate article
// @Tags        Articles
// @Description This request is responsible for recording the user's rating of an article.
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       articleId   path string         true "Article ID"
// @Param       requestBody body request.Rating true "Request body"
// @Success 204 {object} nil                       "Successful request"
// @Failure 400 {object} swagger.HttpError "Badly formatted request"
// @Failure 401 {object} swagger.HttpError "Unauthorized access"
// @Failure 403 {object} swagger.HttpError "Access denied"
// @Failure 404 {object} swagger.HttpError "Requested resource not found"
// @Failure 422 {object} swagger.HttpError "Some of the data provided is invalid"
// @Failure 500 {object} swagger.HttpError "An unexpected error occurred while processing the request"
// @Failure 503 {object} swagger.HttpError "Some of the services/resources are temporarily unavailable"
// @Router /articles/{articleId}/rating [PUT]
func (instance Article) SaveArticleRating(context echo.Context) error {
	userId := utils.GetUserIdFromAuthorizationHeader(context)

	articleIdParameter := context.Param("articleId")
	parameter, parameterDescription := "articleId", "Article ID"
	articleId, httpError := utils.ConvertFromStringToUuid(articleIdParameter, parameter, parameterDescription)
	if httpError != nil {
		log.Warn("Error converting the articleId parameter: ", httpError.Message)
		return context.JSON(httpError.Code, httpError)
	}

	var rating request.Rating
	err := context.Bind(&rating)
	if err != nil {
		log.Warn("Error assigning data from article rating request to DTO: ", err.Error())
		return context.JSON(http.StatusBadRequest, response.NewBadRequestError())
	}

	if rating.Rating != nil && (*rating.Rating < 1 || *rating.Rating > 5) {
		errorMessage := "The article rating value is invalid"
		log.Warnf("%s: (Value: %d; User: %s)", errorMessage, rating.Rating, userId)
		return context.JSON(http.StatusUnprocessableEntity, response.NewHttpError(http.StatusUnprocessableEntity,
			errorMessage))
	}

	err = instance.articleService.SaveArticleRating(userId, articleId, rating.Rating)
	if err != nil {
		if strings.Contains(err.Error(), "user_article_article_fk") {
			log.Warnf("Article %s rated by user %s could not be found: %s", articleId, userId, err.Error())
			return context.JSON(http.StatusNotFound, response.NewHttpError(http.StatusNotFound,
				"Article not found"))
		} else if strings.Contains(err.Error(), "duplicate key") {
			log.Warnf("Could not rate the article %s with user %s: %s", articleId, userId, err.Error())
			return context.JSON(http.StatusForbidden, response.NewForbiddenError())
		} else if strings.Contains(err.Error(), "connection refused") {
			log.Error("Database unavailable: ", err.Error())
			return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
		}

		log.Errorf("Error rating the article %s with user %s: %s", articleId, userId, err.Error())
		return context.JSON(http.StatusInternalServerError, response.NewInternalServerError())
	}

	return context.NoContent(http.StatusNoContent)
}

// SaveArticleToViewLater
// @ID          SaveArticleToViewLater
// @Summary     Add or remove an article from the list of articles bookmarked for later viewing by the user
// @Tags        Articles
// @Description This request is responsible for adding or removing the article from the list of articles bookmarked for later viewing by the user.
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       articleId   path string            true "Article ID"
// @Param       requestBody body request.ViewLater true "Request body"
// @Success 204 {object} nil                       "Successful request"
// @Failure 400 {object} swagger.HttpError "Badly formatted request"
// @Failure 401 {object} swagger.HttpError "Unauthorized access"
// @Failure 403 {object} swagger.HttpError "Access denied"
// @Failure 404 {object} swagger.HttpError "Requested resource not found"
// @Failure 422 {object} swagger.HttpError "Some of the data provided is invalid"
// @Failure 500 {object} swagger.HttpError "An unexpected error occurred while processing the request"
// @Failure 503 {object} swagger.HttpError "Some of the services/resources are temporarily unavailable"
// @Router /articles/{articleId}/view-later [PUT]
func (instance Article) SaveArticleToViewLater(context echo.Context) error {
	userId := utils.GetUserIdFromAuthorizationHeader(context)

	articleIdParameter := context.Param("articleId")
	parameter, parameterDescription := "articleId", "Article ID"
	articleId, httpError := utils.ConvertFromStringToUuid(articleIdParameter, parameter, parameterDescription)
	if httpError != nil {
		log.Warn("Error converting the articleId parameter: ", httpError.Message)
		return context.JSON(httpError.Code, httpError)
	}

	var viewLater request.ViewLater
	err := context.Bind(&viewLater)
	if err != nil {
		log.Warn("Error assigning data from article bookmarking request to view later to DTO: ", err.Error())
		return context.JSON(http.StatusBadRequest, response.NewBadRequestError())
	}

	err = instance.articleService.SaveArticleToViewLater(userId, articleId, viewLater.ViewLater)
	if err != nil {
		if strings.Contains(err.Error(), "user_article_article_fk") {
			log.Warnf("Could not find article %s bookmarked/unbookmarked for later viewing by user %s: %s",
				articleId, userId, err.Error())
			return context.JSON(http.StatusNotFound, response.NewHttpError(http.StatusNotFound,
				"Article not found"))
		} else if strings.Contains(err.Error(), "duplicate key") {
			log.Warnf("Could not update the bookmarking of article %s with user %s: %s",
				articleId, userId, err.Error())
			return context.JSON(http.StatusForbidden, response.NewForbiddenError())
		} else if strings.Contains(err.Error(), "connection refused") {
			log.Error("Database unavailable: ", err.Error())
			return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
		}

		log.Errorf("Error updating the later viewing bookmark of article %s with user %s: %s",
			articleId, userId, err.Error())
		return context.JSON(http.StatusInternalServerError, response.NewInternalServerError())
	}

	return context.NoContent(http.StatusNoContent)
}
