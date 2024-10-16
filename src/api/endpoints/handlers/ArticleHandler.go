package handlers

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"strings"
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
	newsletterService  services.Newsletter
}

func NewArticleHandler(articleService services.Article, resourceService services.Resources,
	propositionService services.Proposition, newsletterService services.Newsletter) *Article {
	return &Article{
		articleService:     articleService,
		resourceService:    resourceService,
		propositionService: propositionService,
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
// @Param       typeId           query string false "Article type ID"
// @Param       content          query string false "Part of the content of the articles, in the title or content"
// @Param       deputyId         query string false "ID of the deputy who drafted the proposition"
// @Param       partyId          query string false "ID of the party that drafted the proposition"
// @Param       externalAuthorId query string false "ID of the external author who drafted the proposition"
// @Param       startDate        query string false "Date from which the articles may have been created. Accepted format: YYYY-MM-DD"
// @Param       endDate          query string false "Date until which the articles may have been created. Accepted format: YYYY-MM-DD"
// @Param       page             query int    false "Page number. By default, it is 1"
// @Param       itemsPerPage     query int    false "Number of articles returned per page. The default is 15 and the allowed values are between 1 and 100"
// @Success 200 {object} response.SwaggerArticlePagination "Successful request"
// @Failure 400 {object} response.SwaggerHttpError         "Badly formatted request"
// @Failure 401 {object} response.SwaggerHttpError         "Unauthorized access"
// @Failure 422 {object} response.SwaggerHttpError         "Some of the data provided is invalid"
// @Failure 500 {object} response.SwaggerHttpError         "An unexpected error occurred while processing the request"
// @Failure 503 {object} response.SwaggerHttpError         "Some of the services/resources are temporarily unavailable"
// @Router /articles [GET]
func (instance Article) GetArticles(context echo.Context) error {
	userId := utils.GetUserIdFromAuthorizationHeader(context)

	var articleFilter filters.ArticleFilter

	typeIdParam := context.QueryParam("typeId")
	if typeIdParam != "" {
		typeId, httpError := utils.ConvertFromStringToUuid(typeIdParam, "Article type ID")
		if httpError != nil {
			log.Error("Error converting the typeId parameter: ", httpError.Message)
			return context.JSON(httpError.Code, httpError)
		}
		articleFilter.TypeId = &typeId
	}

	articleFilter.Content = context.QueryParam("content")

	deputyIdParam := context.QueryParam("deputyId")
	if deputyIdParam != "" {
		deputyId, httpError := utils.ConvertFromStringToUuid(deputyIdParam, "Deputy ID")
		if httpError != nil {
			log.Error("Error converting the deputyId parameter: ", httpError.Message)
			return context.JSON(httpError.Code, httpError)
		}
		articleFilter.DeputyId = &deputyId
	}

	partyIdParam := context.QueryParam("partyId")
	if partyIdParam != "" {
		partyId, httpError := utils.ConvertFromStringToUuid(partyIdParam, "Party ID")
		if httpError != nil {
			log.Error("Error converting the partyId parameter: ", httpError.Message)
			return context.JSON(httpError.Code, httpError)
		}
		articleFilter.PartyId = &partyId
	}

	externalAuthorIdParam := context.QueryParam("externalAuthorId")
	if externalAuthorIdParam != "" {
		externalAuthorId, httpError := utils.ConvertFromStringToUuid(externalAuthorIdParam, "External author ID")
		if httpError != nil {
			log.Error("Error converting the externalAuthorId parameter: ", httpError.Message)
			return context.JSON(httpError.Code, httpError)
		}
		articleFilter.ExternalAuthorId = &externalAuthorId
	}

	startDateParam := context.QueryParam("startDate")
	if startDateParam != "" {
		startDate, httpError := utils.ConvertFromStringToTime(startDateParam, "Start date")
		if httpError != nil {
			log.Error("Error converting the startDate parameter: ", httpError.Message)
			return context.JSON(httpError.Code, httpError)
		}
		articleFilter.StartDate = &startDate
	}

	endDateParam := context.QueryParam("endDate")
	if endDateParam != "" {
		endDate, httpError := utils.ConvertFromStringToTime(endDateParam, "End date")
		if httpError != nil {
			log.Error("Error converting the endDate parameter: ", httpError.Message)
			return context.JSON(httpError.Code, httpError)
		}
		articleFilter.EndDate = &endDate
	}

	if articleFilter.StartDate != nil && articleFilter.EndDate != nil && articleFilter.StartDate.After(*articleFilter.EndDate) {
		errorMessage := fmt.Sprintf("Invalid parameters: The startDate parameter cannot be greater than the endDate parameter")
		log.Error(errorMessage)
		return context.JSON(http.StatusUnprocessableEntity, response.NewHttpError(http.StatusUnprocessableEntity, errorMessage))
	}

	pageParam := context.QueryParam("page")
	if pageParam != "" {
		page, httpError := utils.ConvertFromStringToInt(pageParam, "Page")
		if httpError != nil {
			log.Error("Error converting the page parameter: ", httpError.Message)
			return context.JSON(httpError.Code, httpError)
		}
		articleFilter.PaginationFilter.Page = &page
	}

	itemsPerPageParam := context.QueryParam("itemsPerPage")
	if itemsPerPageParam != "" {
		itemsPerPage, httpError := utils.ConvertFromStringToInt(itemsPerPageParam, "Items per page")
		if httpError != nil {
			log.Error("Error converting the itemsPerPage parameter: ", httpError.Message)
			return context.JSON(httpError.Code, httpError)
		}

		if itemsPerPage > 100 {
			errorMessage := fmt.Sprintf("Invalid parameter: Items per page")
			log.Errorf("Badly formatted request: %s (Value: %d)", errorMessage, itemsPerPage)
			return context.JSON(http.StatusBadRequest, response.NewHttpError(http.StatusBadRequest, errorMessage))
		}

		articleFilter.PaginationFilter.ItemsPerPage = &itemsPerPage
	}

	articleSlice, totalNumberOfArticles, err := instance.articleService.GetArticles(articleFilter, userId)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			log.Error("Database unavailable: ", err.Error())
			return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
		}

		log.Error("Error fetching articles: ", err.Error())
		return context.JSON(http.StatusInternalServerError, response.NewInternalServerError())
	}

	articles := []response.Article{}
	for _, articleData := range articleSlice {
		articles = append(articles, *response.NewArticle(articleData))
	}

	requestResult := response.Pagination{
		Page:         articleFilter.PaginationFilter.GetPage(),
		ItemsPerPage: articleFilter.PaginationFilter.GetItemsPerPage(),
		Total:        totalNumberOfArticles,
		Data:         articles,
	}

	return context.JSON(http.StatusOK, requestResult)
}

// GetTrendingArticles
// @ID          GetTrendingArticles
// @Summary     List trending articles
// @Tags        Articles
// @Description This request is responsible for listing trending articles on the platform.
// @Security    BearerAuth
// @Produce     json
// @Param       typeId           query string false "Article type ID"
// @Param       content          query string false "Part of the content of the articles, in the title or content"
// @Param       deputyId         query string false "ID of the deputy who drafted the proposition"
// @Param       partyId          query string false "ID of the party that drafted the proposition"
// @Param       externalAuthorId query string false "ID of the external author who drafted the proposition"
// @Param       startDate        query string false "Date from which the articles may have been created. Accepted format: YYYY-MM-DD"
// @Param       endDate          query string false "Date until which the articles may have been created. Accepted format: YYYY-MM-DD"
// @Param       page             query int    false "Page number. By default, it is 1"
// @Param       itemsPerPage     query int    false "Number of articles returned per page. The default is 15 and the allowed values are between 1 and 100"
// @Success 200 {object} response.SwaggerArticlePagination "Successful request"
// @Failure 400 {object} response.SwaggerHttpError         "Badly formatted request"
// @Failure 401 {object} response.SwaggerHttpError         "Unauthorized access"
// @Failure 422 {object} response.SwaggerHttpError         "Some of the data provided is invalid"
// @Failure 500 {object} response.SwaggerHttpError         "An unexpected error occurred while processing the request"
// @Failure 503 {object} response.SwaggerHttpError         "Some of the services/resources are temporarily unavailable"
// @Router /articles/trending [GET]
func (instance Article) GetTrendingArticles(context echo.Context) error {
	userId := utils.GetUserIdFromAuthorizationHeader(context)

	var articleFilter filters.ArticleFilter

	typeIdParam := context.QueryParam("typeId")
	if typeIdParam != "" {
		typeId, httpError := utils.ConvertFromStringToUuid(typeIdParam, "Article type ID")
		if httpError != nil {
			log.Error("Error converting the typeId parameter: ", httpError.Message)
			return context.JSON(httpError.Code, httpError)
		}
		articleFilter.TypeId = &typeId
	}

	articleFilter.Content = context.QueryParam("content")

	deputyIdParam := context.QueryParam("deputyId")
	if deputyIdParam != "" {
		deputyId, httpError := utils.ConvertFromStringToUuid(deputyIdParam, "Deputy ID")
		if httpError != nil {
			log.Error("Error converting the deputyId parameter: ", httpError.Message)
			return context.JSON(httpError.Code, httpError)
		}
		articleFilter.DeputyId = &deputyId
	}

	partyIdParam := context.QueryParam("partyId")
	if partyIdParam != "" {
		partyId, httpError := utils.ConvertFromStringToUuid(partyIdParam, "Party ID")
		if httpError != nil {
			log.Error("Error converting the partyId parameter: ", httpError.Message)
			return context.JSON(httpError.Code, httpError)
		}
		articleFilter.PartyId = &partyId
	}

	externalAuthorIdParam := context.QueryParam("externalAuthorId")
	if externalAuthorIdParam != "" {
		externalAuthorId, httpError := utils.ConvertFromStringToUuid(externalAuthorIdParam, "External author ID")
		if httpError != nil {
			log.Error("Error converting the externalAuthorId parameter: ", httpError.Message)
			return context.JSON(httpError.Code, httpError)
		}
		articleFilter.ExternalAuthorId = &externalAuthorId
	}

	startDateParam := context.QueryParam("startDate")
	if startDateParam != "" {
		startDate, httpError := utils.ConvertFromStringToTime(startDateParam, "Start date")
		if httpError != nil {
			log.Error("Error converting the startDate parameter: ", httpError.Message)
			return context.JSON(httpError.Code, httpError)
		}
		articleFilter.StartDate = &startDate
	}

	endDateParam := context.QueryParam("endDate")
	if endDateParam != "" {
		endDate, httpError := utils.ConvertFromStringToTime(endDateParam, "End date")
		if httpError != nil {
			log.Error("Error converting the endDate parameter: ", httpError.Message)
			return context.JSON(httpError.Code, httpError)
		}
		articleFilter.EndDate = &endDate
	}

	if articleFilter.StartDate != nil && articleFilter.EndDate != nil && articleFilter.StartDate.After(*articleFilter.EndDate) {
		errorMessage := fmt.Sprintf("Invalid parameters: The startDate parameter cannot be greater than the endDate parameter")
		log.Error(errorMessage)
		return context.JSON(http.StatusUnprocessableEntity, response.NewHttpError(http.StatusUnprocessableEntity, errorMessage))
	}

	pageParam := context.QueryParam("page")
	if pageParam != "" {
		page, httpError := utils.ConvertFromStringToInt(pageParam, "Page")
		if httpError != nil {
			log.Error("Error converting the page parameter: ", httpError.Message)
			return context.JSON(httpError.Code, httpError)
		}
		articleFilter.PaginationFilter.Page = &page
	}

	itemsPerPageParam := context.QueryParam("itemsPerPage")
	if itemsPerPageParam != "" {
		itemsPerPage, httpError := utils.ConvertFromStringToInt(itemsPerPageParam, "Items per page")
		if httpError != nil {
			log.Error("Error converting the itemsPerPage parameter: ", httpError.Message)
			return context.JSON(httpError.Code, httpError)
		}

		if itemsPerPage > 100 {
			errorMessage := fmt.Sprintf("Invalid parameter: Items per page")
			log.Errorf("Badly formatted request: %s (Value: %d)", errorMessage, itemsPerPage)
			return context.JSON(http.StatusBadRequest, response.NewHttpError(http.StatusBadRequest, errorMessage))
		}

		articleFilter.PaginationFilter.ItemsPerPage = &itemsPerPage
	}

	articleSlice, totalNumberOfArticles, err := instance.articleService.GetTrendingArticles(articleFilter, userId)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			log.Error("Database unavailable: ", err.Error())
			return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
		}

		log.Error("Error fetching trending articles: ", err.Error())
		return context.JSON(http.StatusInternalServerError, response.NewInternalServerError())
	}

	articles := []response.Article{}
	for _, articleData := range articleSlice {
		articles = append(articles, *response.NewArticle(articleData))
	}

	requestResult := response.Pagination{
		Page:         articleFilter.PaginationFilter.GetPage(),
		ItemsPerPage: articleFilter.PaginationFilter.GetItemsPerPage(),
		Total:        totalNumberOfArticles,
		Data:         articles,
	}

	return context.JSON(http.StatusOK, requestResult)
}

// GetTrendingArticlesByTypeId
// @ID          GetTrendingArticlesByTypeId
// @Summary     List trending articles by article types
// @Tags        Articles
// @Description This request is responsible for listing the trending articles by article types.
// @Security    BearerAuth
// @Produce     json
// @Param       articleTypeIds query string false "List of IDs of the types of articles that should be returned (separated by commas). By default, it returns all types"
// @Param       itemsPerType   query int    false "Number of articles returned by type. The default is 5 and the allowed values are between 1 and 20"
// @Success 200 {object} response.SwaggerArticleTypeWithArticles "Successful request"
// @Failure 400 {object} response.SwaggerHttpError               "Badly formatted request"
// @Failure 401 {object} response.SwaggerHttpError               "Unauthorized access"
// @Failure 500 {object} response.SwaggerHttpError               "An unexpected error occurred while processing the request"
// @Failure 503 {object} response.SwaggerHttpError               "Some of the services/resources are temporarily unavailable"
// @Router /articles/trending/type [GET]
func (instance Article) GetTrendingArticlesByTypeId(context echo.Context) error {
	userId := utils.GetUserIdFromAuthorizationHeader(context)

	var articleTypeIds []uuid.UUID
	idsParam := context.QueryParam("articleTypeIds")
	if idsParam != "" {
		articleTypeIdsAsStringSlice := strings.Split(idsParam, ",")
		for index, articleTypeIdAsString := range articleTypeIdsAsStringSlice {
			articleTypeId, httpError := utils.ConvertFromStringToUuid(articleTypeIdAsString,
				fmt.Sprintf("%d°th article type ID", index+1))
			if httpError != nil {
				log.Error("Error converting the articleTypeIds parameter: ", httpError.Message)
				return context.JSON(httpError.Code, httpError)
			}
			articleTypeIds = append(articleTypeIds, articleTypeId)
		}
	}

	itemsPerType := 5
	itemsPerTypeParam := context.QueryParam("itemsPerType")
	if itemsPerTypeParam != "" {
		var httpError *response.HttpError
		itemsPerType, httpError = utils.ConvertFromStringToInt(itemsPerTypeParam, "Items per type")
		if httpError != nil {
			log.Error("Error converting the itemsPerType parameter: ", httpError.Message)
			return context.JSON(httpError.Code, httpError)
		}

		if itemsPerType > 20 {
			errorMessage := fmt.Sprintf("Invalid parameter: Items per type")
			log.Errorf("Badly formatted request: %s (Value: %d)", errorMessage, itemsPerType)
			return context.JSON(http.StatusBadRequest, response.NewHttpError(http.StatusBadRequest, errorMessage))
		}
	}

	articleTypes, err := instance.resourceService.GetArticleTypes(articleTypeIds)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			log.Error("Database unavailable: ", err.Error())
			return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
		}

		log.Error("Error fetching article type details: ", err.Error())
		return context.JSON(http.StatusInternalServerError, response.NewInternalServerError())
	}

	var articleTypeSlice []response.ArticleType
	for _, articleType := range articleTypes {
		articles, err := instance.articleService.GetTrendingArticlesByTypeId(articleType.Id(), itemsPerType, userId)
		if err != nil {
			if strings.Contains(err.Error(), "connection refused") {
				log.Error("Database unavailable: ", err.Error())
				return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
			}

			log.Errorf("Error fetching articles of article type %s: %s", articleType.Id(), err.Error())
			return context.JSON(http.StatusInternalServerError, response.NewInternalServerError())
		}

		articleTypeResponse := response.NewArticleType(articleType)

		var articleSlice []response.Article
		for _, articleData := range articles {
			articleSlice = append(articleSlice, *response.NewArticle(articleData))
		}

		articleTypeResponse.Articles = articleSlice

		articleTypeSlice = append(articleTypeSlice, *articleTypeResponse)
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
// @Param       typeId           query string false "Article type ID"
// @Param       content          query string false "Part of the content of the articles, in the title or content"
// @Param       deputyId         query string false "ID of the deputy who drafted the proposition"
// @Param       partyId          query string false "ID of the party that drafted the proposition"
// @Param       externalAuthorId query string false "ID of the external author who drafted the proposition"
// @Param       startDate        query string false "Date from which the articles may have been created. Accepted format: YYYY-MM-DD"
// @Param       endDate          query string false "Date until which the articles may have been created. Accepted format: YYYY-MM-DD"
// @Param       page             query int    false "Page number. By default, it is 1"
// @Param       itemsPerPage     query int    false "Number of articles returned per page. The default is 15 and the allowed values are between 1 and 100"
// @Success 200 {object} response.SwaggerArticlePagination "Successful request"
// @Failure 400 {object} response.SwaggerHttpError         "Badly formatted request"
// @Failure 401 {object} response.SwaggerHttpError         "Unauthorized access"
// @Failure 422 {object} response.SwaggerHttpError         "Some of the data provided is invalid"
// @Failure 500 {object} response.SwaggerHttpError         "An unexpected error occurred while processing the request"
// @Failure 503 {object} response.SwaggerHttpError         "Some of the services/resources are temporarily unavailable"
// @Router /articles/view-later [GET]
func (instance Article) GetArticlesToViewLater(context echo.Context) error {
	userId := utils.GetUserIdFromAuthorizationHeader(context)

	var articleFilter filters.ArticleFilter

	typeIdParam := context.QueryParam("typeId")
	if typeIdParam != "" {
		typeId, httpError := utils.ConvertFromStringToUuid(typeIdParam, "Article type ID")
		if httpError != nil {
			log.Error("Error converting the typeId parameter: ", httpError.Message)
			return context.JSON(httpError.Code, httpError)
		}
		articleFilter.TypeId = &typeId
	}

	articleFilter.Content = context.QueryParam("content")

	deputyIdParam := context.QueryParam("deputyId")
	if deputyIdParam != "" {
		deputyId, httpError := utils.ConvertFromStringToUuid(deputyIdParam, "Deputy ID")
		if httpError != nil {
			log.Error("Error converting the deputyId parameter: ", httpError.Message)
			return context.JSON(httpError.Code, httpError)
		}
		articleFilter.DeputyId = &deputyId
	}

	partyIdParam := context.QueryParam("partyId")
	if partyIdParam != "" {
		partyId, httpError := utils.ConvertFromStringToUuid(partyIdParam, "Party ID")
		if httpError != nil {
			log.Error("Error converting the partyId parameter: ", httpError.Message)
			return context.JSON(httpError.Code, httpError)
		}
		articleFilter.PartyId = &partyId
	}

	externalAuthorIdParam := context.QueryParam("externalAuthorId")
	if externalAuthorIdParam != "" {
		externalAuthorId, httpError := utils.ConvertFromStringToUuid(externalAuthorIdParam, "External author ID")
		if httpError != nil {
			log.Error("Error converting the externalAuthorId parameter: ", httpError.Message)
			return context.JSON(httpError.Code, httpError)
		}
		articleFilter.ExternalAuthorId = &externalAuthorId
	}

	startDateParam := context.QueryParam("startDate")
	if startDateParam != "" {
		startDate, httpError := utils.ConvertFromStringToTime(startDateParam, "Start date")
		if httpError != nil {
			log.Error("Error converting the startDate parameter: ", httpError.Message)
			return context.JSON(httpError.Code, httpError)
		}
		articleFilter.StartDate = &startDate
	}

	endDateParam := context.QueryParam("endDate")
	if endDateParam != "" {
		endDate, httpError := utils.ConvertFromStringToTime(endDateParam, "End date")
		if httpError != nil {
			log.Error("Error converting the endDate parameter: ", httpError.Message)
			return context.JSON(httpError.Code, httpError)
		}
		articleFilter.EndDate = &endDate
	}

	if articleFilter.StartDate != nil && articleFilter.EndDate != nil && articleFilter.StartDate.After(*articleFilter.EndDate) {
		errorMessage := fmt.Sprintf("Invalid parameters: The startDate parameter cannot be greater than the endDate parameter")
		log.Error(errorMessage)
		return context.JSON(http.StatusUnprocessableEntity, response.NewHttpError(http.StatusUnprocessableEntity, errorMessage))
	}

	pageParam := context.QueryParam("page")
	if pageParam != "" {
		page, httpError := utils.ConvertFromStringToInt(pageParam, "Page")
		if httpError != nil {
			log.Error("Error converting the page parameter: ", httpError.Message)
			return context.JSON(httpError.Code, httpError)
		}
		articleFilter.PaginationFilter.Page = &page
	}

	itemsPerPageParam := context.QueryParam("itemsPerPage")
	if itemsPerPageParam != "" {
		itemsPerPage, httpError := utils.ConvertFromStringToInt(itemsPerPageParam, "Items per page")
		if httpError != nil {
			log.Error("Error converting the itemsPerPage parameter: ", httpError.Message)
			return context.JSON(httpError.Code, httpError)
		}

		if itemsPerPage > 100 {
			errorMessage := fmt.Sprintf("Invalid parameter: Items per page")
			log.Errorf("Badly formatted request: %s (Value: %d)", errorMessage, itemsPerPage)
			return context.JSON(http.StatusBadRequest, response.NewHttpError(http.StatusBadRequest, errorMessage))
		}

		articleFilter.PaginationFilter.ItemsPerPage = &itemsPerPage
	}

	articleSlice, totalNumberOfArticles, err := instance.articleService.GetArticlesToViewLater(articleFilter, userId)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			log.Error("Database unavailable: ", err.Error())
			return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
		}

		log.Errorf("Error fetching articles bookmarked for later viewing by user %s: %s", userId, err.Error())
		return context.JSON(http.StatusInternalServerError, response.NewInternalServerError())
	}

	articles := []response.Article{}
	for _, articleData := range articleSlice {
		articles = append(articles, *response.NewArticle(articleData))
	}

	requestResult := response.Pagination{
		Page:         articleFilter.PaginationFilter.GetPage(),
		ItemsPerPage: articleFilter.PaginationFilter.GetItemsPerPage(),
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
// @Success 200 {object} response.SwaggerPropositionArticle "Successful request"
// @Failure 400 {object} response.SwaggerHttpError          "Badly formatted request"
// @Failure 401 {object} response.SwaggerHttpError          "Unauthorized access"
// @Failure 404 {object} response.SwaggerHttpError          "Requested resource not found"
// @Failure 422 {object} response.SwaggerHttpError          "Some of the data provided is invalid"
// @Failure 500 {object} response.SwaggerHttpError          "An unexpected error occurred while processing the request"
// @Failure 503 {object} response.SwaggerHttpError          "Some of the services/resources are temporarily unavailable"
// @Router /articles/{articleId}/proposition [GET]
func (instance Article) GetPropositionArticleById(context echo.Context) error {
	userId := utils.GetUserIdFromAuthorizationHeader(context)

	articleId, httpError := utils.ConvertFromStringToUuid(context.Param("articleId"), "Article ID")
	if httpError != nil {
		log.Error("Error converting the articleId parameter: ", httpError.Message)
		return context.JSON(httpError.Code, httpError)
	}

	propositionData, err := instance.propositionService.GetPropositionByArticleId(articleId, userId)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			log.Errorf("Proposition article %s could not be found: %s", articleId, err.Error())
			return context.JSON(http.StatusNotFound, response.NewHttpError(http.StatusNotFound,
				"Proposition article not found"))
		} else if strings.Contains(err.Error(), "connection refused") {
			log.Error("Database unavailable: ", err.Error())
			return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
		}

		log.Errorf("Error fetching proposition article %s: %s", articleId, err.Error())
		return context.JSON(http.StatusInternalServerError, response.NewInternalServerError())
	}

	propositionArticle := response.NewPropositionArticle(*propositionData)

	newsletterArticle, err := instance.articleService.GetNewsletterArticleByPropositionId(propositionData.Id(), userId)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			log.Error("Database unavailable: ", err.Error())
			return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
		} else if !strings.Contains(err.Error(), "no rows") {
			log.Errorf("Error fetching newsletter article of proposition %s (Article: %s): %s",
				propositionData.Id(), articleId, err.Error())
			return context.JSON(http.StatusInternalServerError, response.NewInternalServerError())
		}
	} else {
		propositionArticle.NewsletterArticle = response.NewArticle(*newsletterArticle)
	}

	return context.JSON(http.StatusOK, propositionArticle)
}

// GetNewsletterArticleById
// @ID          GetNewsletterArticleById
// @Summary     Get article details by ID (Only for newsletter articles)
// @Tags        Articles
// @Description This request is responsible for looking up the details of an article of a newsletter by the article ID.
// @Security    BearerAuth
// @Produce     json
// @Param       articleId path string true "Article ID"
// @Success 200 {object} response.SwaggerNewsletterArticle "Successful request"
// @Failure 400 {object} response.SwaggerHttpError         "Badly formatted request"
// @Failure 401 {object} response.SwaggerHttpError         "Unauthorized access"
// @Failure 404 {object} response.SwaggerHttpError         "Requested resource not found"
// @Failure 422 {object} response.SwaggerHttpError         "Some of the data provided is invalid"
// @Failure 500 {object} response.SwaggerHttpError         "An unexpected error occurred while processing the request"
// @Failure 503 {object} response.SwaggerHttpError         "Some of the services/resources are temporarily unavailable"
// @Router /articles/{articleId}/newsletter [GET]
func (instance Article) GetNewsletterArticleById(context echo.Context) error {
	userId := utils.GetUserIdFromAuthorizationHeader(context)

	articleId, httpError := utils.ConvertFromStringToUuid(context.Param("articleId"), "Article ID")
	if httpError != nil {
		log.Error("Error converting the articleId parameter: ", httpError.Message)
		return context.JSON(httpError.Code, httpError)
	}

	newsletterData, err := instance.newsletterService.GetNewsletterByArticleId(articleId, userId)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			log.Errorf("Newsletter article %s could not be found: %s", articleId, err.Error())
			return context.JSON(http.StatusNotFound, response.NewHttpError(http.StatusNotFound,
				"Newsletter article not found"))
		} else if strings.Contains(err.Error(), "connection refused") {
			log.Error("Database unavailable: ", err.Error())
			return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
		}

		log.Errorf("Error fetching newsletter article %s: %s", articleId, err.Error())
		return context.JSON(http.StatusInternalServerError, response.NewInternalServerError())
	}

	newsletterArticle := response.NewNewsletterArticle(*newsletterData)

	propositionArticles, err := instance.articleService.GetPropositionArticlesByNewsletterId(newsletterData.Id(), userId)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			log.Error("Database unavailable: ", err.Error())
			return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
		}

		log.Errorf("Error fetching proposition articles of newsletter %s (Article: %s): %s",
			newsletterData.Id(), articleId, err.Error())
		return context.JSON(http.StatusInternalServerError, response.NewInternalServerError())
	}

	for _, articleData := range propositionArticles {
		newsletterArticle.PropositionArticles = append(newsletterArticle.PropositionArticles, *response.NewArticle(articleData))
	}

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
// @Param       articleId path string         true "Article ID"
// @Param       body      body request.Rating true "Request body"
// @Success 204 {object} nil                       "Successful request"
// @Failure 400 {object} response.SwaggerHttpError "Badly formatted request"
// @Failure 401 {object} response.SwaggerHttpError "Unauthorized access"
// @Failure 403 {object} response.SwaggerHttpError "Access denied"
// @Failure 404 {object} response.SwaggerHttpError "Requested resource not found"
// @Failure 422 {object} response.SwaggerHttpError "Some of the data provided is invalid"
// @Failure 500 {object} response.SwaggerHttpError "An unexpected error occurred while processing the request"
// @Failure 503 {object} response.SwaggerHttpError "Some of the services/resources are temporarily unavailable"
// @Router /articles/{articleId}/rating [PUT]
func (instance Article) SaveArticleRating(context echo.Context) error {
	userId := utils.GetUserIdFromAuthorizationHeader(context)

	articleId, httpError := utils.ConvertFromStringToUuid(context.Param("articleId"), "Article ID")
	if httpError != nil {
		log.Error("Error converting the articleId parameter: ", httpError.Message)
		return context.JSON(httpError.Code, httpError)
	}

	var rating request.Rating
	err := context.Bind(&rating)
	if err != nil {
		log.Error("Error assigning data from article rating request to DTO: ", err.Error())
		return context.JSON(http.StatusBadRequest, response.NewBadRequestError())
	}

	if rating.Rating < 1 || rating.Rating > 5 {
		errorMessage := "The article rating value is invalid"
		log.Errorf("%s: (Value: %d; User: %s)", errorMessage, rating.Rating, userId)
		return context.JSON(http.StatusUnprocessableEntity, response.NewHttpError(http.StatusUnprocessableEntity, errorMessage))
	}

	err = instance.articleService.SaveArticleRating(userId, articleId, rating.Rating)
	if err != nil {
		if strings.Contains(err.Error(), "user_article_article_fk") {
			log.Errorf("Article %s rated by user %s could not be found: %s", articleId, userId, err.Error())
			return context.JSON(http.StatusNotFound, response.NewHttpError(http.StatusNotFound,
				"Article not found"))
		} else if strings.Contains(err.Error(), "duplicate key") {
			log.Errorf("Could not rate the article %s with user %s: %s", articleId, userId, err.Error())
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
// @Param       articleId path string            true "Article ID"
// @Param       body      body request.ViewLater true "Request body"
// @Success 204 {object} nil                       "Successful request"
// @Failure 400 {object} response.SwaggerHttpError "Badly formatted request"
// @Failure 401 {object} response.SwaggerHttpError "Unauthorized access"
// @Failure 403 {object} response.SwaggerHttpError "Access denied"
// @Failure 404 {object} response.SwaggerHttpError "Requested resource not found"
// @Failure 422 {object} response.SwaggerHttpError "Some of the data provided is invalid"
// @Failure 500 {object} response.SwaggerHttpError "An unexpected error occurred while processing the request"
// @Failure 503 {object} response.SwaggerHttpError "Some of the services/resources are temporarily unavailable"
// @Router /articles/{articleId}/view-later [PUT]
func (instance Article) SaveArticleToViewLater(context echo.Context) error {
	userId := utils.GetUserIdFromAuthorizationHeader(context)

	articleId, httpError := utils.ConvertFromStringToUuid(context.Param("articleId"), "Article ID")
	if httpError != nil {
		log.Error("Error converting the articleId parameter: ", httpError.Message)
		return context.JSON(httpError.Code, httpError)
	}

	var viewLater request.ViewLater
	err := context.Bind(&viewLater)
	if err != nil {
		log.Error("Error assigning data from article bookmarking request to view later to DTO: ",
			err.Error())
		return context.JSON(http.StatusBadRequest, response.NewBadRequestError())
	}

	err = instance.articleService.SaveArticleToViewLater(userId, articleId, viewLater.ViewLater)
	if err != nil {
		if strings.Contains(err.Error(), "user_article_article_fk") {
			log.Errorf("Could not find article %s bookmarked/unbookmarked for later viewing by user %s: %s",
				articleId, userId, err.Error())
			return context.JSON(http.StatusNotFound, response.NewHttpError(http.StatusNotFound,
				"Article not found"))
		} else if strings.Contains(err.Error(), "duplicate key") {
			log.Errorf("Could not update the bookmarking of article %s with user %s: %s",
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
