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
// @Summary     Listagem das matérias mais recentes
// @Tags        Matérias
// @Description Esta requisição é responsável por retornar as matérias mais recentes disponíveis na plataforma Você na Câmara.
// @Security    BearerAuth
// @Produce     json
// @Param       content          query string false "Parte do conteúdo das matérias, no título ou conteúdo."
// @Param       deputyId         query string false "ID do deputado que elaborou a proposição."
// @Param       partyId          query string false "ID do partido que elaborou a proposição."
// @Param       externalAuthorId query string false "ID do autor externo que elaborou a proposição."
// @Param       startDate        query string false "Data inicial para submissão da proposta. Formato aceito: YYYY-MM-DD"
// @Param       endDate          query string false "Data final para submissão da proposta. Formato aceito: YYYY-MM-DD"
// @Param       type             query string false "Tipo das matérias. Por padrão retorna todos os tipos. Valores permitidos: 'Proposição', 'Boletim'."
// @Param       page             query int    false "Número da página. Por padrão é 1."
// @Param       itemsPerPage     query int    false "Quantidade de matérias retornadas por página. Por padrão é 15 e os valores permitidos são entre 1 e 100."
// @Success 200 {object} response.SwaggerArticlePagination "Requisição realizada com sucesso."
// @Failure 400 {object} response.SwaggerHttpError         "Requisição mal formulada."
// @Failure 422 {object} response.SwaggerHttpError         "Requisição não processada devido a algum dos dados enviados serem inválidos."
// @Failure 500 {object} response.SwaggerHttpError         "Ocorreu um erro inesperado durante o processamento da requisição."
// @Failure 503 {object} response.SwaggerHttpError         "Algum dos serviços/recursos está temporariamente indisponível."
// @Router /articles [GET]
func (instance Article) GetArticles(context echo.Context) error {
	userId := utils.GetUserIdFromAuthorizationHeader(context)

	var articleFilter filters.ArticleFilter

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

	typeParam := context.QueryParam("type")
	if typeParam != "" {
		if typeParam != "Proposição" && typeParam != "Boletim" {
			errorMessage := fmt.Sprintf("Parâmetro inválido: type")
			log.Errorf("Requisição mal formulada: %s (Valor: %s)", errorMessage, typeParam)
			return context.JSON(http.StatusUnprocessableEntity, response.NewHttpError(http.StatusUnprocessableEntity, errorMessage))
		} else if typeParam == "Boletim" && (articleFilter.DeputyId != nil || articleFilter.PartyId != nil ||
			articleFilter.ExternalAuthorId != nil) {
			errorMessage := fmt.Sprintf("Parâmetros inválidos: O parâmetro type quando igual a 'Boletim' não " +
				"permite o uso dos demais parâmetros query deputyId, partyId e externalAuthorId.")
			log.Error(errorMessage)
			return context.JSON(http.StatusUnprocessableEntity, response.NewHttpError(http.StatusUnprocessableEntity, errorMessage))
		}
		articleFilter.Type = typeParam
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

	articleList, totalNumberOfArticles, err := instance.articleService.GetArticles(articleFilter, userId)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			log.Error("Banco de dados indisponível: ", err.Error())
			return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
		}

		log.Error("Erro ao buscar as matérias: ", err.Error())
		return context.JSON(http.StatusInternalServerError, response.NewInternalServerError())
	}

	var articles []response.Article
	for _, articleData := range articleList {
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
// @Summary     Listagem das matérias em alta
// @Tags        Matérias
// @Description Esta requisição é responsável por retornar as matérias em alta disponíveis na plataforma Você na Câmara.
// @Security    BearerAuth
// @Produce     json
// @Param       content          query string false "Parte do conteúdo das matérias, no título ou conteúdo."
// @Param       deputyId         query string false "ID do deputado que elaborou a proposição."
// @Param       partyId          query string false "ID do partido que elaborou a proposição."
// @Param       externalAuthorId query string false "ID do autor externo que elaborou a proposição."
// @Param       startDate        query string false "Data inicial para submissão da proposta. Formato aceito: YYYY-MM-DD"
// @Param       endDate          query string false "Data final para submissão da proposta. Formato aceito: YYYY-MM-DD"
// @Param       type             query string false "Tipo das matérias. Por padrão retorna todos os tipos. Valores permitidos: 'Proposição', 'Boletim'."
// @Param       page             query int    false "Número da página. Por padrão é 1."
// @Param       itemsPerPage     query int    false "Quantidade de matérias retornadas por página. Por padrão é 15 e os valores permitidos são entre 1 e 100."
// @Success 200 {object} response.SwaggerArticlePagination "Requisição realizada com sucesso."
// @Failure 400 {object} response.SwaggerHttpError         "Requisição mal formulada."
// @Failure 422 {object} response.SwaggerHttpError         "Requisição não processada devido a algum dos dados enviados serem inválidos."
// @Failure 500 {object} response.SwaggerHttpError         "Ocorreu um erro inesperado durante o processamento da requisição."
// @Failure 503 {object} response.SwaggerHttpError         "Algum dos serviços/recursos está temporariamente indisponível."
// @Router /articles/trending [GET]
func (instance Article) GetTrendingArticles(context echo.Context) error {
	userId := utils.GetUserIdFromAuthorizationHeader(context)

	var articleFilter filters.ArticleFilter

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

	typeParam := context.QueryParam("type")
	if typeParam != "" {
		if typeParam != "Proposição" && typeParam != "Boletim" {
			errorMessage := fmt.Sprintf("Parâmetro inválido: type")
			log.Errorf("Requisição mal formulada: %s (Valor: %s)", errorMessage, typeParam)
			return context.JSON(http.StatusUnprocessableEntity, response.NewHttpError(http.StatusUnprocessableEntity, errorMessage))
		} else if typeParam == "Boletim" && (articleFilter.DeputyId != nil || articleFilter.PartyId != nil ||
			articleFilter.ExternalAuthorId != nil) {
			errorMessage := fmt.Sprintf("Parâmetros inválidos: O parâmetro type quando igual a 'Boletim' não " +
				"permite o uso dos demais parâmetros query deputyId, partyId e externalAuthorId.")
			log.Error(errorMessage)
			return context.JSON(http.StatusUnprocessableEntity, response.NewHttpError(http.StatusUnprocessableEntity, errorMessage))
		}
		articleFilter.Type = typeParam
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

	articleList, totalNumberOfArticles, err := instance.articleService.GetTrendingArticles(articleFilter, userId)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			log.Error("Banco de dados indisponível: ", err.Error())
			return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
		}

		log.Error("Erro ao buscar as matérias em alta: ", err.Error())
		return context.JSON(http.StatusInternalServerError, response.NewInternalServerError())
	}

	var articles []response.Article
	for _, articleData := range articleList {
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

// GetTrendingArticlesByPropositionType
// @ID          GetTrendingArticlesByPropositionType
// @Summary     Listagem das matérias em alta pelos tipos das proposições
// @Tags        Matérias
// @Description Esta requisição é responsável por retornar as matérias em alta pelos tipos das proposições.
// @Security    BearerAuth
// @Produce     json
// @Param       propositionTypeIds query string false "Lista com os IDs dos tipos das proposições que devem ser retornados (Separados por vírgula). Por padrão retorna todos."
// @Param       itemsPerType       query int    false "Quantidade de matérias retornadas por tipo. Por padrão é 5 e os valores permitidos são entre 1 e 20."
// @Success 200 {object} response.SwaggerPropositionTypeWithPropositions "Requisição realizada com sucesso."
// @Failure 400 {object} response.SwaggerHttpError                       "Requisição mal formulada."
// @Failure 500 {object} response.SwaggerHttpError                       "Ocorreu um erro inesperado durante o processamento da requisição."
// @Failure 503 {object} response.SwaggerHttpError                       "Algum dos serviços/recursos está temporariamente indisponível."
// @Router /articles/trending/proposition-type [GET]
func (instance Article) GetTrendingArticlesByPropositionType(context echo.Context) error {
	userId := utils.GetUserIdFromAuthorizationHeader(context)

	var propositionTypeIds []uuid.UUID
	propositionTypeIdsParam := context.QueryParam("propositionTypeIds")
	if propositionTypeIdsParam != "" {
		propositionTypeIdsAsStringSlice := strings.Split(propositionTypeIdsParam, ",")
		for index, propositionTypeIdAsString := range propositionTypeIdsAsStringSlice {
			propositionTypeId, httpError := utils.ConvertFromStringToUuid(propositionTypeIdAsString,
				fmt.Sprintf("Id do %d° tipo de proposição", index))
			if httpError != nil {
				log.Error("Erro ao converter o parâmetro propositionTypeIds: ", httpError.Message)
				return context.JSON(httpError.Code, httpError)
			}
			propositionTypeIds = append(propositionTypeIds, propositionTypeId)
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

	propositionTypes, err := instance.resourceService.GetPropositionTypes(propositionTypeIds)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			log.Error("Banco de dados indisponível: ", err.Error())
			return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
		}

		log.Error("Erro ao buscar os detalhes dos tipos de proposições no banco de dados: ", err.Error())
		return context.JSON(http.StatusInternalServerError, response.NewInternalServerError())
	}

	var propositionTypeList []response.PropositionType
	for _, propositionType := range propositionTypes {
		articles, err := instance.articleService.GetTrendingArticlesByPropositionType(
			propositionType.Id(), itemsPerType, userId)
		if err != nil {
			if strings.Contains(err.Error(), "connection refused") {
				log.Error("Banco de dados indisponível: ", err.Error())
				return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
			}

			log.Errorf("Erro ao buscar as matérias do tipo de proposição %s: %s", propositionType.Id(), err.Error())
			return context.JSON(http.StatusInternalServerError, response.NewInternalServerError())
		}

		propositionTypeResponse := response.NewPropositionType(propositionType)

		var propositionArticles []response.Article
		for _, articleData := range articles {
			propositionArticles = append(propositionArticles, *response.NewArticle(articleData))
		}

		propositionTypeResponse.PropositionsArticles = propositionArticles

		propositionTypeList = append(propositionTypeList, *propositionTypeResponse)
	}

	return context.JSON(http.StatusOK, propositionTypeList)
}

// GetArticlesToViewLater
// @ID          GetArticlesToViewLater
// @Summary     Listagem das matérias marcadas para ver depois
// @Tags        Matérias
// @Description Esta requisição é responsável por retornar as matérias marcadas para ver depois na plataforma Você na Câmara.
// @Security    BearerAuth
// @Produce     json
// @Param       content          query string false "Parte do conteúdo das matérias, no título ou conteúdo."
// @Param       deputyId         query string false "ID do deputado que elaborou a proposição."
// @Param       partyId          query string false "ID do partido que elaborou a proposição."
// @Param       externalAuthorId query string false "ID do autor externo que elaborou a proposição."
// @Param       startDate        query string false "Data inicial para submissão da proposta. Formato aceito: YYYY-MM-DD"
// @Param       endDate          query string false "Data final para submissão da proposta. Formato aceito: YYYY-MM-DD"
// @Param       type             query string false "Tipo das matérias. Por padrão retorna todos os tipos. Valores permitidos: 'Proposição', 'Boletim'."
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

	typeParam := context.QueryParam("type")
	if typeParam != "" {
		if typeParam != "Proposição" && typeParam != "Boletim" {
			errorMessage := fmt.Sprintf("Parâmetro inválido: type")
			log.Errorf("Requisição mal formulada: %s (Valor: %s)", errorMessage, typeParam)
			return context.JSON(http.StatusUnprocessableEntity, response.NewHttpError(http.StatusUnprocessableEntity, errorMessage))
		} else if typeParam == "Boletim" && (articleFilter.DeputyId != nil || articleFilter.PartyId != nil ||
			articleFilter.ExternalAuthorId != nil) {
			errorMessage := fmt.Sprintf("Parâmetros inválidos: O parâmetro type quando igual a 'Boletim' não " +
				"permite o uso dos demais parâmetros query deputyId, partyId e externalAuthorId.")
			log.Error(errorMessage)
			return context.JSON(http.StatusUnprocessableEntity, response.NewHttpError(http.StatusUnprocessableEntity, errorMessage))
		}
		articleFilter.Type = typeParam
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

	articleList, totalNumberOfArticles, err := instance.articleService.GetArticlesToViewLater(articleFilter, userId)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			log.Error("Banco de dados indisponível: ", err.Error())
			return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
		}

		log.Error("Erro ao buscar as matérias marcadas para ver depois pelo usuário %s: %s", userId, err.Error())
		return context.JSON(http.StatusInternalServerError, response.NewInternalServerError())
	}

	var articles []response.Article
	for _, articleData := range articleList {
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
// @Summary     Busca dos detalhes da matéria pelo ID do tipo proposição
// @Tags        Matérias
// @Description Esta requisição é responsável por retornar os detalhes da matéria pelo ID do tipo proposição.
// @Security    BearerAuth
// @Produce     json
// @Param       articleId path string true "ID da matéria"
// @Success 200 {object} response.SwaggerPropositionArticle "Requisição realizada com sucesso."
// @Failure 400 {object} response.SwaggerHttpError          "Requisição mal formulada."
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
// @Summary     Busca dos detalhes da matéria pelo ID do tipo boletim
// @Tags        Matérias
// @Description Esta requisição é responsável por retornar os detalhes da matéria pelo ID do tipo boletim.
// @Security    BearerAuth
// @Produce     json
// @Param       articleId path string true "ID da matéria"
// @Success 200 {object} response.SwaggerNewsletterArticle "Requisição realizada com sucesso."
// @Failure 400 {object} response.SwaggerHttpError         "Requisição mal formulada."
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
// @Description Esta requisição é responsável pelo registro da avaliação do usuário sobre uma matéria.
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       articleId path string         true "ID da matéria"
// @Param       body      body request.Rating true "JSON com todos os dados necessários para que o login seja realizado."
// @Success 204 {object} nil                       "Requisição realizada com sucesso."
// @Failure 400 {object} response.SwaggerHttpError "Requisição mal formulada."
// @Failure 401 {object} response.SwaggerHttpError "Acesso não autorizado."
// @Failure 403 {object} response.SwaggerHttpError "Acesso negado."
// @Failure 404 {object} response.SwaggerHttpError "Recurso solicitado não encontrado."
// @Failure 422 {object} response.SwaggerHttpError "Requisição não processada devido a algum dos dados enviados serem inválidos."
// @Failure 500 {object} response.SwaggerHttpError "Ocorreu um erro inesperado durante o processamento da requisição."
// @Failure 503 {object} response.SwaggerHttpError "Algum dos serviços/recursos está temporariamente indisponível."
// @Router /articles/{articleId}/rating [PATCH]
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
		log.Error("Erro ao atribuir os dados da requisição ao DTO: ", err.Error())
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
// @Summary     Salvar matéria para ver depois
// @Tags        Matérias
// @Description Esta requisição é responsável por salvar a matéria na lista de matérias marcadas para ver depois do usuário.
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       articleId path string            true "ID da matéria"
// @Param       body      body request.ViewLater true "JSON com todos os dados necessários para que o login seja realizado."
// @Success 204 {object} nil                       "Requisição realizada com sucesso."
// @Failure 400 {object} response.SwaggerHttpError "Requisição mal formulada."
// @Failure 401 {object} response.SwaggerHttpError "Acesso não autorizado."
// @Failure 403 {object} response.SwaggerHttpError "Acesso negado."
// @Failure 404 {object} response.SwaggerHttpError "Recurso solicitado não encontrado."
// @Failure 422 {object} response.SwaggerHttpError "Requisição não processada devido a algum dos dados enviados serem inválidos."
// @Failure 500 {object} response.SwaggerHttpError "Ocorreu um erro inesperado durante o processamento da requisição."
// @Failure 503 {object} response.SwaggerHttpError "Algum dos serviços/recursos está temporariamente indisponível."
// @Router /articles/{articleId}/view-later [PATCH]
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
		log.Error("Erro ao atribuir os dados da requisição ao DTO: ", err.Error())
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
