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
// @Summary     Listar matérias mais recentes
// @Tags        Matérias
// @Description Esta requisição é responsável por listar as matérias mais recentes disponíveis na plataforma.
// @Security    BearerAuth
// @Produce     json
// @Param       typeId           query string false "ID do tipo da matéria."
// @Param       content          query string false "Parte do conteúdo das matérias, no título ou conteúdo."
// @Param       deputyId         query string false "ID do deputado que elaborou a proposição."
// @Param       partyId          query string false "ID do partido que elaborou a proposição."
// @Param       externalAuthorId query string false "ID do autor externo que elaborou a proposição."
// @Param       startDate        query string false "Data a partir da qual as matérias podem ter sido criadas. Formato aceito: YYYY-MM-DD"
// @Param       endDate          query string false "Data até a qual as matérias podem ter sido criadas. Formato aceito: YYYY-MM-DD"
// @Param       page             query int    false "Número da página. Por padrão é 1."
// @Param       itemsPerPage     query int    false "Quantidade de matérias retornadas por página. Por padrão é 15 e os valores permitidos são entre 1 e 100."
// @Success 200 {object} response.SwaggerArticlePagination "Requisição realizada com sucesso."
// @Failure 400 {object} response.SwaggerHttpError         "Requisição mal formulada."
// @Failure 401 {object} response.SwaggerHttpError         "Acesso não autorizado."
// @Failure 422 {object} response.SwaggerHttpError         "Requisição não processada devido a algum dos dados enviados serem inválidos."
// @Failure 500 {object} response.SwaggerHttpError         "Ocorreu um erro inesperado durante o processamento da requisição."
// @Failure 503 {object} response.SwaggerHttpError         "Algum dos serviços/recursos está temporariamente indisponível."
// @Router /articles [GET]
func (instance Article) GetArticles(context echo.Context) error {
	userId := utils.GetUserIdFromAuthorizationHeader(context)

	var articleFilter filters.ArticleFilter

	typeIdParam := context.QueryParam("typeId")
	if typeIdParam != "" {
		typeId, httpError := utils.ConvertFromStringToUuid(typeIdParam, "ID do tipo da matéria")
		if httpError != nil {
			log.Error("Erro ao converter o parâmetro typeId: ", httpError.Message)
			return context.JSON(httpError.Code, httpError)
		}
		articleFilter.TypeId = &typeId
	}

	articleFilter.Content = context.QueryParam("content")

	deputyIdParam := context.QueryParam("deputyId")
	if deputyIdParam != "" {
		deputyId, httpError := utils.ConvertFromStringToUuid(deputyIdParam, "ID do deputado")
		if httpError != nil {
			log.Error("Erro ao converter o parâmetro deputyId: ", httpError.Message)
			return context.JSON(httpError.Code, httpError)
		}
		articleFilter.DeputyId = &deputyId
	}

	partyIdParam := context.QueryParam("partyId")
	if partyIdParam != "" {
		partyId, httpError := utils.ConvertFromStringToUuid(partyIdParam, "ID do partido")
		if httpError != nil {
			log.Error("Erro ao converter o parâmetro partyId: ", httpError.Message)
			return context.JSON(httpError.Code, httpError)
		}
		articleFilter.PartyId = &partyId
	}

	externalAuthorIdParam := context.QueryParam("externalAuthorId")
	if externalAuthorIdParam != "" {
		externalAuthorId, httpError := utils.ConvertFromStringToUuid(externalAuthorIdParam, "ID do autor externo")
		if httpError != nil {
			log.Error("Erro ao converter o parâmetro externalAuthorId: ", httpError.Message)
			return context.JSON(httpError.Code, httpError)
		}
		articleFilter.ExternalAuthorId = &externalAuthorId
	}

	startDateParam := context.QueryParam("startDate")
	if startDateParam != "" {
		startDate, httpError := utils.ConvertFromStringToTime(startDateParam, "Data inicial")
		if httpError != nil {
			log.Error("Erro ao converter o parâmetro startDate: ", httpError.Message)
			return context.JSON(httpError.Code, httpError)
		}
		articleFilter.StartDate = &startDate
	}

	endDateParam := context.QueryParam("endDate")
	if endDateParam != "" {
		endDate, httpError := utils.ConvertFromStringToTime(endDateParam, "Data final")
		if httpError != nil {
			log.Error("Erro ao converter o parâmetro endDate: ", httpError.Message)
			return context.JSON(httpError.Code, httpError)
		}
		articleFilter.EndDate = &endDate
	}

	if articleFilter.StartDate != nil && articleFilter.EndDate != nil && articleFilter.StartDate.After(*articleFilter.EndDate) {
		errorMessage := fmt.Sprintf("Parâmetros inválidos: O parâmetro startDate não pode ser maior que o parâmetro endDate")
		log.Error(errorMessage)
		return context.JSON(http.StatusUnprocessableEntity, response.NewHttpError(http.StatusUnprocessableEntity, errorMessage))
	}

	pageParam := context.QueryParam("page")
	if pageParam != "" {
		page, httpError := utils.ConvertFromStringToInt(pageParam, "Página")
		if httpError != nil {
			log.Error("Erro ao converter o parâmetro page: ", httpError.Message)
			return context.JSON(httpError.Code, httpError)
		}
		articleFilter.PaginationFilter.Page = &page
	}

	itemsPerPageParam := context.QueryParam("itemsPerPage")
	if itemsPerPageParam != "" {
		itemsPerPage, httpError := utils.ConvertFromStringToInt(itemsPerPageParam, "Itens por página")
		if httpError != nil {
			log.Error("Erro ao converter o parâmetro itemsPerPage: ", httpError.Message)
			return context.JSON(httpError.Code, httpError)
		}

		if itemsPerPage > 100 {
			errorMessage := fmt.Sprintf("Parâmetro inválido: Itens por página")
			log.Errorf("Requisição mal formulada: %s (Valor: %s)", errorMessage, itemsPerPage)
			return context.JSON(http.StatusBadRequest, response.NewHttpError(http.StatusBadRequest, errorMessage))
		}

		articleFilter.PaginationFilter.ItemsPerPage = &itemsPerPage
	}

	articleSlice, totalNumberOfArticles, err := instance.articleService.GetArticles(articleFilter, userId)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			log.Error("Banco de dados indisponível: ", err.Error())
			return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
		}

		log.Error("Erro ao buscar as matérias: ", err.Error())
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
// @Summary     Listar matérias em alta
// @Tags        Matérias
// @Description Esta requisição é responsável por listar as matérias em alta disponíveis na plataforma.
// @Security    BearerAuth
// @Produce     json
// @Param       typeId           query string false "ID do tipo da matéria."
// @Param       content          query string false "Parte do conteúdo das matérias, no título ou conteúdo."
// @Param       deputyId         query string false "ID do deputado que elaborou a proposição."
// @Param       partyId          query string false "ID do partido que elaborou a proposição."
// @Param       externalAuthorId query string false "ID do autor externo que elaborou a proposição."
// @Param       startDate        query string false "Data a partir da qual as matérias podem ter sido criadas. Formato aceito: YYYY-MM-DD"
// @Param       endDate          query string false "Data até a qual as matérias podem ter sido criadas. Formato aceito: YYYY-MM-DD"
// @Param       page             query int    false "Número da página. Por padrão é 1."
// @Param       itemsPerPage     query int    false "Quantidade de matérias retornadas por página. Por padrão é 15 e os valores permitidos são entre 1 e 100."
// @Success 200 {object} response.SwaggerArticlePagination "Requisição realizada com sucesso."
// @Failure 400 {object} response.SwaggerHttpError         "Requisição mal formulada."
// @Failure 401 {object} response.SwaggerHttpError         "Acesso não autorizado."
// @Failure 422 {object} response.SwaggerHttpError         "Requisição não processada devido a algum dos dados enviados serem inválidos."
// @Failure 500 {object} response.SwaggerHttpError         "Ocorreu um erro inesperado durante o processamento da requisição."
// @Failure 503 {object} response.SwaggerHttpError         "Algum dos serviços/recursos está temporariamente indisponível."
// @Router /articles/trending [GET]
func (instance Article) GetTrendingArticles(context echo.Context) error {
	userId := utils.GetUserIdFromAuthorizationHeader(context)

	var articleFilter filters.ArticleFilter

	typeIdParam := context.QueryParam("typeId")
	if typeIdParam != "" {
		typeId, httpError := utils.ConvertFromStringToUuid(typeIdParam, "ID do tipo da matéria")
		if httpError != nil {
			log.Error("Erro ao converter o parâmetro typeId: ", httpError.Message)
			return context.JSON(httpError.Code, httpError)
		}
		articleFilter.TypeId = &typeId
	}

	articleFilter.Content = context.QueryParam("content")

	deputyIdParam := context.QueryParam("deputyId")
	if deputyIdParam != "" {
		deputyId, httpError := utils.ConvertFromStringToUuid(deputyIdParam, "ID do deputado")
		if httpError != nil {
			log.Error("Erro ao converter o parâmetro deputyId: ", httpError.Message)
			return context.JSON(httpError.Code, httpError)
		}
		articleFilter.DeputyId = &deputyId
	}

	partyIdParam := context.QueryParam("partyId")
	if partyIdParam != "" {
		partyId, httpError := utils.ConvertFromStringToUuid(partyIdParam, "ID do partido")
		if httpError != nil {
			log.Error("Erro ao converter o parâmetro partyId: ", httpError.Message)
			return context.JSON(httpError.Code, httpError)
		}
		articleFilter.PartyId = &partyId
	}

	externalAuthorIdParam := context.QueryParam("externalAuthorId")
	if externalAuthorIdParam != "" {
		externalAuthorId, httpError := utils.ConvertFromStringToUuid(externalAuthorIdParam, "ID do autor externo")
		if httpError != nil {
			log.Error("Erro ao converter o parâmetro externalAuthorId: ", httpError.Message)
			return context.JSON(httpError.Code, httpError)
		}
		articleFilter.ExternalAuthorId = &externalAuthorId
	}

	startDateParam := context.QueryParam("startDate")
	if startDateParam != "" {
		startDate, httpError := utils.ConvertFromStringToTime(startDateParam, "Data inicial")
		if httpError != nil {
			log.Error("Erro ao converter o parâmetro startDate: ", httpError.Message)
			return context.JSON(httpError.Code, httpError)
		}
		articleFilter.StartDate = &startDate
	}

	endDateParam := context.QueryParam("endDate")
	if endDateParam != "" {
		endDate, httpError := utils.ConvertFromStringToTime(endDateParam, "Data final")
		if httpError != nil {
			log.Error("Erro ao converter o parâmetro endDate: ", httpError.Message)
			return context.JSON(httpError.Code, httpError)
		}
		articleFilter.EndDate = &endDate
	}

	if articleFilter.StartDate != nil && articleFilter.EndDate != nil && articleFilter.StartDate.After(*articleFilter.EndDate) {
		errorMessage := fmt.Sprintf("Parâmetros inválidos: O parâmetro startDate não pode ser maior que o parâmetro endDate")
		log.Error(errorMessage)
		return context.JSON(http.StatusUnprocessableEntity, response.NewHttpError(http.StatusUnprocessableEntity, errorMessage))
	}

	pageParam := context.QueryParam("page")
	if pageParam != "" {
		page, httpError := utils.ConvertFromStringToInt(pageParam, "Página")
		if httpError != nil {
			log.Error("Erro ao converter o parâmetro page: ", httpError.Message)
			return context.JSON(httpError.Code, httpError)
		}
		articleFilter.PaginationFilter.Page = &page
	}

	itemsPerPageParam := context.QueryParam("itemsPerPage")
	if itemsPerPageParam != "" {
		itemsPerPage, httpError := utils.ConvertFromStringToInt(itemsPerPageParam, "Itens por página")
		if httpError != nil {
			log.Error("Erro ao converter o parâmetro itemsPerPage: ", httpError.Message)
			return context.JSON(httpError.Code, httpError)
		}

		if itemsPerPage > 100 {
			errorMessage := fmt.Sprintf("Parâmetro inválido: Itens por página")
			log.Errorf("Requisição mal formulada: %s (Valor: %s)", errorMessage, itemsPerPage)
			return context.JSON(http.StatusBadRequest, response.NewHttpError(http.StatusBadRequest, errorMessage))
		}

		articleFilter.PaginationFilter.ItemsPerPage = &itemsPerPage
	}

	articleSlice, totalNumberOfArticles, err := instance.articleService.GetTrendingArticles(articleFilter, userId)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			log.Error("Banco de dados indisponível: ", err.Error())
			return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
		}

		log.Error("Erro ao buscar as matérias em alta: ", err.Error())
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
// @Summary     Listar matérias em alta pelos tipos de matérias
// @Tags        Matérias
// @Description Esta requisição é responsável por listar as matérias em alta pelos tipos de matérias.
// @Security    BearerAuth
// @Produce     json
// @Param       ids          query string false "Lista com os IDs dos tipos de matérias que devem ser retornados (Separados por vírgula). Por padrão retorna todos."
// @Param       itemsPerType query int    false "Quantidade de matérias retornadas por tipo. Por padrão é 5 e os valores permitidos são entre 1 e 20."
// @Success 200 {object} response.SwaggerArticleTypeWithArticles "Requisição realizada com sucesso."
// @Failure 400 {object} response.SwaggerHttpError               "Requisição mal formulada."
// @Failure 401 {object} response.SwaggerHttpError               "Acesso não autorizado."
// @Failure 500 {object} response.SwaggerHttpError               "Ocorreu um erro inesperado durante o processamento da requisição."
// @Failure 503 {object} response.SwaggerHttpError               "Algum dos serviços/recursos está temporariamente indisponível."
// @Router /articles/trending/type [GET]
func (instance Article) GetTrendingArticlesByTypeId(context echo.Context) error {
	userId := utils.GetUserIdFromAuthorizationHeader(context)

	var articleTypeIds []uuid.UUID
	idsParam := context.QueryParam("ids")
	if idsParam != "" {
		articleTypeIdsAsStringSlice := strings.Split(idsParam, ",")
		for index, articleTypeIdAsString := range articleTypeIdsAsStringSlice {
			articleTypeId, httpError := utils.ConvertFromStringToUuid(articleTypeIdAsString,
				fmt.Sprintf("Id do %d° tipo de matéria", index+1))
			if httpError != nil {
				log.Error("Erro ao converter o parâmetro ids: ", httpError.Message)
				return context.JSON(httpError.Code, httpError)
			}
			articleTypeIds = append(articleTypeIds, articleTypeId)
		}
	}

	itemsPerType := 5
	itemsPerTypeParam := context.QueryParam("itemsPerType")
	if itemsPerTypeParam != "" {
		var httpError *response.HttpError
		itemsPerType, httpError = utils.ConvertFromStringToInt(itemsPerTypeParam, "Itens por tipo")
		if httpError != nil {
			log.Error("Erro ao converter o parâmetro itemsPerType: ", httpError.Message)
			return context.JSON(httpError.Code, httpError)
		}

		if itemsPerType > 20 {
			errorMessage := fmt.Sprintf("Parâmetro inválido: Itens por tipo")
			log.Errorf("Requisição mal formulada: %s (Valor: %s)", errorMessage, itemsPerType)
			return context.JSON(http.StatusBadRequest, response.NewHttpError(http.StatusBadRequest, errorMessage))
		}
	}

	articleTypes, err := instance.resourceService.GetArticleTypes(articleTypeIds)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			log.Error("Banco de dados indisponível: ", err.Error())
			return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
		}

		log.Error("Erro ao buscar os detalhes dos tipos de matéria no banco de dados: ", err.Error())
		return context.JSON(http.StatusInternalServerError, response.NewInternalServerError())
	}

	var articleTypeSlice []response.ArticleType
	for _, articleType := range articleTypes {
		articles, err := instance.articleService.GetTrendingArticlesByTypeId(articleType.Id(), itemsPerType, userId)
		if err != nil {
			if strings.Contains(err.Error(), "connection refused") {
				log.Error("Banco de dados indisponível: ", err.Error())
				return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
			}

			log.Errorf("Erro ao buscar as matérias do tipo de matéria %s: %s", articleType.Id(), err.Error())
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
// @Summary     Listar matérias marcadas para ver depois pelo usuário
// @Tags        Matérias
// @Description Esta requisição é responsável por listar as matérias marcadas para ver depois pelo usuário na plataforma. As matérias serão listadas na ordem que o usuário realizou a marcação das matérias.
// @Security    BearerAuth
// @Produce     json
// @Param       typeId           query string false "ID do tipo da matéria."
// @Param       content          query string false "Parte do conteúdo das matérias, no título ou conteúdo."
// @Param       deputyId         query string false "ID do deputado que elaborou a proposição."
// @Param       partyId          query string false "ID do partido que elaborou a proposição."
// @Param       externalAuthorId query string false "ID do autor externo que elaborou a proposição."
// @Param       startDate        query string false "Data a partir da qual as matérias podem ter sido criadas. Formato aceito: YYYY-MM-DD"
// @Param       endDate          query string false "Data até a qual as matérias podem ter sido criadas. Formato aceito: YYYY-MM-DD"
// @Param       page             query int    false "Número da página. Por padrão é 1."
// @Param       itemsPerPage     query int    false "Quantidade de matérias retornadas por página. Por padrão é 15 e os valores permitidos são entre 1 e 100."
// @Success 200 {object} response.SwaggerArticlePagination "Requisição realizada com sucesso."
// @Failure 400 {object} response.SwaggerHttpError         "Requisição mal formulada."
// @Failure 401 {object} response.SwaggerHttpError         "Acesso não autorizado."
// @Failure 422 {object} response.SwaggerHttpError         "Requisição não processada devido a algum dos dados enviados serem inválidos."
// @Failure 500 {object} response.SwaggerHttpError         "Ocorreu um erro inesperado durante o processamento da requisição."
// @Failure 503 {object} response.SwaggerHttpError         "Algum dos serviços/recursos está temporariamente indisponível."
// @Router /articles/view-later [GET]
func (instance Article) GetArticlesToViewLater(context echo.Context) error {
	userId := utils.GetUserIdFromAuthorizationHeader(context)

	var articleFilter filters.ArticleFilter

	typeIdParam := context.QueryParam("typeId")
	if typeIdParam != "" {
		typeId, httpError := utils.ConvertFromStringToUuid(typeIdParam, "ID do tipo da matéria")
		if httpError != nil {
			log.Error("Erro ao converter o parâmetro typeId: ", httpError.Message)
			return context.JSON(httpError.Code, httpError)
		}
		articleFilter.TypeId = &typeId
	}

	articleFilter.Content = context.QueryParam("content")

	deputyIdParam := context.QueryParam("deputyId")
	if deputyIdParam != "" {
		deputyId, httpError := utils.ConvertFromStringToUuid(deputyIdParam, "ID do deputado")
		if httpError != nil {
			log.Error("Erro ao converter o parâmetro deputyId: ", httpError.Message)
			return context.JSON(httpError.Code, httpError)
		}
		articleFilter.DeputyId = &deputyId
	}

	partyIdParam := context.QueryParam("partyId")
	if partyIdParam != "" {
		partyId, httpError := utils.ConvertFromStringToUuid(partyIdParam, "ID do partido")
		if httpError != nil {
			log.Error("Erro ao converter o parâmetro partyId: ", httpError.Message)
			return context.JSON(httpError.Code, httpError)
		}
		articleFilter.PartyId = &partyId
	}

	externalAuthorIdParam := context.QueryParam("externalAuthorId")
	if externalAuthorIdParam != "" {
		externalAuthorId, httpError := utils.ConvertFromStringToUuid(externalAuthorIdParam, "ID do autor externo")
		if httpError != nil {
			log.Error("Erro ao converter o parâmetro externalAuthorId: ", httpError.Message)
			return context.JSON(httpError.Code, httpError)
		}
		articleFilter.ExternalAuthorId = &externalAuthorId
	}

	startDateParam := context.QueryParam("startDate")
	if startDateParam != "" {
		startDate, httpError := utils.ConvertFromStringToTime(startDateParam, "Data inicial")
		if httpError != nil {
			log.Error("Erro ao converter o parâmetro startDate: ", httpError.Message)
			return context.JSON(httpError.Code, httpError)
		}
		articleFilter.StartDate = &startDate
	}

	endDateParam := context.QueryParam("endDate")
	if endDateParam != "" {
		endDate, httpError := utils.ConvertFromStringToTime(endDateParam, "Data final")
		if httpError != nil {
			log.Error("Erro ao converter o parâmetro endDate: ", httpError.Message)
			return context.JSON(httpError.Code, httpError)
		}
		articleFilter.EndDate = &endDate
	}

	if articleFilter.StartDate != nil && articleFilter.EndDate != nil && articleFilter.StartDate.After(*articleFilter.EndDate) {
		errorMessage := fmt.Sprintf("Parâmetros inválidos: O parâmetro startDate não pode ser maior que o parâmetro endDate")
		log.Error(errorMessage)
		return context.JSON(http.StatusUnprocessableEntity, response.NewHttpError(http.StatusUnprocessableEntity, errorMessage))
	}

	pageParam := context.QueryParam("page")
	if pageParam != "" {
		page, httpError := utils.ConvertFromStringToInt(pageParam, "Página")
		if httpError != nil {
			log.Error("Erro ao converter o parâmetro page: ", httpError.Message)
			return context.JSON(httpError.Code, httpError)
		}
		articleFilter.PaginationFilter.Page = &page
	}

	itemsPerPageParam := context.QueryParam("itemsPerPage")
	if itemsPerPageParam != "" {
		itemsPerPage, httpError := utils.ConvertFromStringToInt(itemsPerPageParam, "Itens por página")
		if httpError != nil {
			log.Error("Erro ao converter o parâmetro itemsPerPage: ", httpError.Message)
			return context.JSON(httpError.Code, httpError)
		}

		if itemsPerPage > 100 {
			errorMessage := fmt.Sprintf("Parâmetro inválido: Itens por página")
			log.Errorf("Requisição mal formulada: %s (Valor: %s)", errorMessage, itemsPerPage)
			return context.JSON(http.StatusBadRequest, response.NewHttpError(http.StatusBadRequest, errorMessage))
		}

		articleFilter.PaginationFilter.ItemsPerPage = &itemsPerPage
	}

	articleSlice, totalNumberOfArticles, err := instance.articleService.GetArticlesToViewLater(articleFilter, userId)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			log.Error("Banco de dados indisponível: ", err.Error())
			return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
		}

		log.Error("Erro ao buscar as matérias marcadas para ver depois pelo usuário %s: %s", userId, err.Error())
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
// @Summary     Buscar detalhes de uma matéria dos tipos de proposições pelo ID
// @Tags        Matérias
// @Description Esta requisição é responsável por buscar os detalhes de uma matéria dos tipos de proposições pelo ID.
// @Security    BearerAuth
// @Produce     json
// @Param       articleId path string true "ID da matéria"
// @Success 200 {object} response.SwaggerPropositionArticle "Requisição realizada com sucesso."
// @Failure 400 {object} response.SwaggerHttpError          "Requisição mal formulada."
// @Failure 401 {object} response.SwaggerHttpError          "Acesso não autorizado."
// @Failure 404 {object} response.SwaggerHttpError          "Recurso solicitado não encontrado."
// @Failure 422 {object} response.SwaggerHttpError          "Requisição não processada devido a algum dos dados enviados serem inválidos."
// @Failure 500 {object} response.SwaggerHttpError          "Ocorreu um erro inesperado durante o processamento da requisição."
// @Failure 503 {object} response.SwaggerHttpError          "Algum dos serviços/recursos está temporariamente indisponível."
// @Router /articles/{articleId}/proposition [GET]
func (instance Article) GetPropositionArticleById(context echo.Context) error {
	userId := utils.GetUserIdFromAuthorizationHeader(context)

	articleId, httpError := utils.ConvertFromStringToUuid(context.Param("articleId"), "ID da matéria")
	if httpError != nil {
		log.Error("Erro ao converter o parâmetro articleId: ", httpError.Message)
		return context.JSON(httpError.Code, httpError)
	}

	propositionData, err := instance.propositionService.GetPropositionByArticleId(articleId, userId)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			log.Errorf("Não foi possível encontrar a matéria da proposição %s: %s", articleId, err.Error())
			return context.JSON(http.StatusNotFound, response.NewHttpError(http.StatusNotFound,
				"Matéria da proposição não encontrada"))
		} else if strings.Contains(err.Error(), "connection refused") {
			log.Error("Banco de dados indisponível: ", err.Error())
			return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
		}

		log.Errorf("Erro ao buscar matéria da proposição %s: %s", articleId, err.Error())
		return context.JSON(http.StatusInternalServerError, response.NewInternalServerError())
	}

	propositionArticle := response.NewPropositionArticle(*propositionData)

	newsletterArticle, err := instance.articleService.GetNewsletterArticleByPropositionId(propositionData.Id(), userId)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			log.Error("Banco de dados indisponível: ", err.Error())
			return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
		} else if !strings.Contains(err.Error(), "no rows") {
			log.Errorf("Erro ao buscar matéria do boletim da proposição %s (Matéria: %s): %s",
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
// @Summary     Buscar detalhes de uma matéria do tipo boletim pelo ID
// @Tags        Matérias
// @Description Esta requisição é responsável por buscar os detalhes de uma matéria do tipo boletim pelo ID.
// @Security    BearerAuth
// @Produce     json
// @Param       articleId path string true "ID da matéria"
// @Success 200 {object} response.SwaggerNewsletterArticle "Requisição realizada com sucesso."
// @Failure 400 {object} response.SwaggerHttpError         "Requisição mal formulada."
// @Failure 401 {object} response.SwaggerHttpError         "Acesso não autorizado."
// @Failure 404 {object} response.SwaggerHttpError         "Recurso solicitado não encontrado."
// @Failure 422 {object} response.SwaggerHttpError         "Requisição não processada devido a algum dos dados enviados serem inválidos."
// @Failure 500 {object} response.SwaggerHttpError         "Ocorreu um erro inesperado durante o processamento da requisição."
// @Failure 503 {object} response.SwaggerHttpError         "Algum dos serviços/recursos está temporariamente indisponível."
// @Router /articles/{articleId}/newsletter [GET]
func (instance Article) GetNewsletterArticleById(context echo.Context) error {
	userId := utils.GetUserIdFromAuthorizationHeader(context)

	articleId, httpError := utils.ConvertFromStringToUuid(context.Param("articleId"), "ID da matéria")
	if httpError != nil {
		log.Error("Erro ao converter o parâmetro articleId: ", httpError.Message)
		return context.JSON(httpError.Code, httpError)
	}

	newsletterData, err := instance.newsletterService.GetNewsletterByArticleId(articleId, userId)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			log.Errorf("Não foi possível encontrar a matéria do boletim %s: %s", articleId, err.Error())
			return context.JSON(http.StatusNotFound, response.NewHttpError(http.StatusNotFound,
				"Matéria do boletim não encontrada"))
		} else if strings.Contains(err.Error(), "connection refused") {
			log.Error("Banco de dados indisponível: ", err.Error())
			return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
		}

		log.Errorf("Erro ao buscar matéria do boletim %s: %s", articleId, err.Error())
		return context.JSON(http.StatusInternalServerError, response.NewInternalServerError())
	}

	newsletterArticle := response.NewNewsletterArticle(*newsletterData)

	propositionArticles, err := instance.articleService.GetPropositionArticlesByNewsletterId(newsletterData.Id(), userId)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			log.Error("Banco de dados indisponível: ", err.Error())
			return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
		}

		log.Errorf("Erro ao buscar matérias das proposições do boletim %s (Matéria: %s): %s",
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
// @Summary     Avaliar matéria
// @Tags        Matérias
// @Description Esta requisição é responsável pelo registro da avaliação de uma matéria pelo usuário.
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       articleId path string         true "ID da matéria"
// @Param       body      body request.Rating true "JSON com todos os dados necessários para que a avaliação da matéria seja realizada."
// @Success 204 {object} nil                       "Requisição realizada com sucesso."
// @Failure 400 {object} response.SwaggerHttpError "Requisição mal formulada."
// @Failure 401 {object} response.SwaggerHttpError "Acesso não autorizado."
// @Failure 403 {object} response.SwaggerHttpError "Acesso negado."
// @Failure 404 {object} response.SwaggerHttpError "Recurso solicitado não encontrado."
// @Failure 422 {object} response.SwaggerHttpError "Requisição não processada devido a algum dos dados enviados serem inválidos."
// @Failure 500 {object} response.SwaggerHttpError "Ocorreu um erro inesperado durante o processamento da requisição."
// @Failure 503 {object} response.SwaggerHttpError "Algum dos serviços/recursos está temporariamente indisponível."
// @Router /articles/{articleId}/rating [PUT]
func (instance Article) SaveArticleRating(context echo.Context) error {
	userId := utils.GetUserIdFromAuthorizationHeader(context)

	articleId, httpError := utils.ConvertFromStringToUuid(context.Param("articleId"), "ID da matéria")
	if httpError != nil {
		log.Error("Erro ao converter o parâmetro articleId: ", httpError.Message)
		return context.JSON(httpError.Code, httpError)
	}

	var rating request.Rating
	err := context.Bind(&rating)
	if err != nil {
		log.Error("Erro ao atribuir os dados da requisição de avaliação da matéria ao DTO: ", err.Error())
		return context.JSON(http.StatusBadRequest, response.NewBadRequestError())
	}

	if rating.Rating < 1 || rating.Rating > 5 {
		errorMessage := "O valor da avaliação da matéria é inválido"
		log.Errorf("%s: (Valor: %d; Usuário: %s)", errorMessage, rating.Rating, userId)
		return context.JSON(http.StatusUnprocessableEntity, response.NewHttpError(http.StatusUnprocessableEntity, errorMessage))
	}

	err = instance.articleService.SaveArticleRating(userId, articleId, rating.Rating)
	if err != nil {
		if strings.Contains(err.Error(), "user_article_article_fk") {
			log.Errorf("Não foi possível encontrar a matéria %s avaliada pelo usuário %s: %s",
				articleId, userId, err.Error())
			return context.JSON(http.StatusNotFound, response.NewHttpError(http.StatusNotFound,
				"Matéria não encontrada"))
		} else if strings.Contains(err.Error(), "duplicate key") {
			log.Errorf("Não foi possível avaliar a matéria %s com o usuário %s: %s", articleId, userId, err.Error())
			return context.JSON(http.StatusForbidden, response.NewForbiddenError())
		} else if strings.Contains(err.Error(), "connection refused") {
			log.Error("Banco de dados indisponível: ", err.Error())
			return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
		}

		log.Errorf("Erro ao avaliar matéria %s com o usuário %s: %s", articleId, userId, err.Error())
		return context.JSON(http.StatusInternalServerError, response.NewInternalServerError())
	}

	return context.NoContent(http.StatusNoContent)
}

// SaveArticleToViewLater
// @ID          SaveArticleToViewLater
// @Summary     Adicionar ou remover matéria da lista de matérias marcadas para ver depois pelo usuário
// @Tags        Matérias
// @Description Esta requisição é responsável por adicionar ou remover a matéria da lista de matérias marcadas para ver depois pelo usuário.
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       articleId path string            true "ID da matéria"
// @Param       body      body request.ViewLater true "JSON com todos os dados necessários para marcar/desmarcar a matéria para ver depois."
// @Success 204 {object} nil                       "Requisição realizada com sucesso."
// @Failure 400 {object} response.SwaggerHttpError "Requisição mal formulada."
// @Failure 401 {object} response.SwaggerHttpError "Acesso não autorizado."
// @Failure 403 {object} response.SwaggerHttpError "Acesso negado."
// @Failure 404 {object} response.SwaggerHttpError "Recurso solicitado não encontrado."
// @Failure 422 {object} response.SwaggerHttpError "Requisição não processada devido a algum dos dados enviados serem inválidos."
// @Failure 500 {object} response.SwaggerHttpError "Ocorreu um erro inesperado durante o processamento da requisição."
// @Failure 503 {object} response.SwaggerHttpError "Algum dos serviços/recursos está temporariamente indisponível."
// @Router /articles/{articleId}/view-later [PUT]
func (instance Article) SaveArticleToViewLater(context echo.Context) error {
	userId := utils.GetUserIdFromAuthorizationHeader(context)

	articleId, httpError := utils.ConvertFromStringToUuid(context.Param("articleId"), "ID da matéria")
	if httpError != nil {
		log.Error("Erro ao converter o parâmetro articleId: ", httpError.Message)
		return context.JSON(httpError.Code, httpError)
	}

	var viewLater request.ViewLater
	err := context.Bind(&viewLater)
	if err != nil {
		log.Error("Erro ao atribuir os dados da requisição de marcação da matéria para ver depois ao DTO: ",
			err.Error())
		return context.JSON(http.StatusBadRequest, response.NewBadRequestError())
	}

	err = instance.articleService.SaveArticleToViewLater(userId, articleId, viewLater.ViewLater)
	if err != nil {
		if strings.Contains(err.Error(), "user_article_article_fk") {
			log.Errorf("Não foi possível encontrar a matéria %s marcada/desmarcada para ver depois pelo usuário "+
				"%s: %s", articleId, userId, err.Error())
			return context.JSON(http.StatusNotFound, response.NewHttpError(http.StatusNotFound,
				"Matéria não encontrada"))
		} else if strings.Contains(err.Error(), "duplicate key") {
			log.Errorf("Não foi possível atualizar a marcação da matéria %s com o usuário %s: %s",
				articleId, userId, err.Error())
			return context.JSON(http.StatusForbidden, response.NewForbiddenError())
		} else if strings.Contains(err.Error(), "connection refused") {
			log.Error("Banco de dados indisponível: ", err.Error())
			return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
		}

		log.Errorf("Erro ao atualizar marcação de ver depois da matéria %s com o usuário %s: %s", articleId,
			userId, err.Error())
		return context.JSON(http.StatusInternalServerError, response.NewInternalServerError())
	}

	return context.NoContent(http.StatusNoContent)
}
