package handlers

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"vnc-read-api/api/endpoints/dto/response"
	"vnc-read-api/core/filters"
	"vnc-read-api/core/interfaces/services"
)

type News struct {
	service services.News
}

func NewNewsHandler(service services.News) *News {
	return &News{service: service}
}

// GetNews
// @ID 			GetNews
// @Summary 	Listagem das matérias mais recentes
// @Tags 		Matérias
// @Description Esta requisição é responsável por retornar as matérias mais recentes disponíveis na plataforma Você na Câmara.
// @Produce		json
// @Param content        query string false "Parte do conteúdo das matérias, no título ou conteúdo."
// @Param deputyId       query string false "ID do deputado que elaborou a proposição."
// @Param partyId        query string false "ID do partido que elaborou a proposição."
// @Param organizationId query string false "ID da organização que elaborou a proposição."
// @Param startDate      query string false "Data inicial para submissão da proposta. Formato aceito: YYYY-MM-DD"
// @Param endDate        query string false "Data final para submissão da proposta. Formato aceito: YYYY-MM-DD"
// @Param type           query string false "Tipo das matérias. Por padrão retorna todos os tipos. Valores permitidos: 'Proposição', 'Boletim'."
// @Param page           query int    false "Número da página. Por padrão é 1."
// @Param itemsPerPage   query int    false "Quantidade de matérias retornadas por página. Por padrão é 15."
// @Success 200 {array}  response.SwaggerNewsPagination "Requisição bem sucedida"
// @Failure 400 {object} response.SwaggerError          "Algum dado informado durante a requisição é inválido"
// @Failure 500 {object} response.SwaggerError          "Ocorreu um erro inesperado durante o processamento da requisição"
// @Router /news [get]
func (instance News) GetNews(context echo.Context) error {
	var newsFilter filters.NewsFilter

	newsFilter.Content = context.QueryParam("content")

	deputyIdParam := context.QueryParam("deputyId")
	if deputyIdParam != "" {
		deputyId, err := convertToUuid(deputyIdParam, "ID do deputado")
		if err != nil {
			log.Error(err)
			return context.JSON(http.StatusBadRequest, response.NewError(err.Error()))
		}
		newsFilter.DeputyId = &deputyId
	}

	partyIdParam := context.QueryParam("partyId")
	if partyIdParam != "" {
		partyId, err := convertToUuid(partyIdParam, "ID do partido")
		if err != nil {
			log.Error(err)
			return context.JSON(http.StatusBadRequest, response.NewError(err.Error()))
		}
		newsFilter.PartyId = &partyId
	}

	organizationIdParam := context.QueryParam("organizationId")
	if organizationIdParam != "" {
		organizationId, err := convertToUuid(organizationIdParam, "ID do organização")
		if err != nil {
			log.Error(err)
			return context.JSON(http.StatusBadRequest, response.NewError(err.Error()))
		}
		newsFilter.OrganizationId = &organizationId
	}

	startDateParam := context.QueryParam("startDate")
	if startDateParam != "" {
		startDate, err := convertToTime(startDateParam, "Data inicial")
		if err != nil {
			log.Error(err)
			return context.JSON(http.StatusBadRequest, response.NewError(err.Error()))
		}
		newsFilter.StartDate = &startDate
	}

	endDateParam := context.QueryParam("endDate")
	if endDateParam != "" {
		endDate, err := convertToTime(endDateParam, "Data final")
		if err != nil {
			log.Error(err)
			return context.JSON(http.StatusBadRequest, response.NewError(err.Error()))
		}
		newsFilter.EndDate = &endDate
	}

	typeParam := context.QueryParam("type")
	if typeParam != "" {
		if typeParam != "Proposição" && typeParam != "Boletim" {
			err := errors.New(fmt.Sprintf("Parâmetro inválido: type"))
			log.Errorf("Requisição mal formulada: %s (Valor: %s)", err, typeParam)
			return context.JSON(http.StatusBadRequest, response.NewError(err.Error()))
		} else if typeParam == "Boletim" && (newsFilter.DeputyId != nil || newsFilter.PartyId != nil || newsFilter.OrganizationId != nil) {
			err := errors.New(fmt.Sprintf("Parâmetro inválido: type igual 'Boletim' não permite o uso dos demais parâmetros query deputyId, partyId e organizationId"))
			log.Error("Parâmetro inválido: type igual 'Boletim' não permite o uso dos demais parâmetros query deputyId, partyId e organizationId")
			return context.JSON(http.StatusBadRequest, response.NewError(err.Error()))
		}
		newsFilter.Type = typeParam
	}

	pageParam := context.QueryParam("page")
	if pageParam != "" {
		page, err := convertToInt(pageParam, "Página")
		if err != nil {
			log.Error(err)
			return context.JSON(http.StatusBadRequest, response.NewError(err.Error()))
		}
		newsFilter.PaginationFilter.Page = &page
	}

	itemsPerPageParam := context.QueryParam("itemsPerPage")
	if itemsPerPageParam != "" {
		perPage, err := convertToInt(itemsPerPageParam, "Itens por página")
		if err != nil {
			log.Error(err)
			return context.JSON(http.StatusBadRequest, response.NewError(err.Error()))
		}
		newsFilter.PaginationFilter.ItemsPerPage = &perPage
	}

	newsList, totalNumberOfNews, err := instance.service.GetNews(newsFilter)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, response.NewError(err.Error()))
	}

	var newsResponse []response.News
	for _, news := range newsList {
		newsResponse = append(newsResponse, *response.NewNews(news))
	}

	requestResult := response.Pagination{
		Page:         newsFilter.PaginationFilter.GetPage(),
		ItensPerPage: newsFilter.PaginationFilter.GetItemsPerPage(),
		Total:        totalNumberOfNews,
		Data:         newsResponse,
	}

	return context.JSON(http.StatusOK, requestResult)
}

// GetTrending
// @ID 			GetTrending
// @Summary 	Listagem das matérias mais populares
// @Tags 		Matérias
// @Description Esta requisição é responsável por retornar as matérias mais populares disponíveis na plataforma Você na Câmara.
// @Produce		json
// @Param content        query string false "Parte do conteúdo das matérias, no título ou conteúdo."
// @Param deputyId       query string false "ID do deputado que elaborou a proposição."
// @Param partyId        query string false "ID do partido que elaborou a proposição."
// @Param organizationId query string false "ID da organização que elaborou a proposição."
// @Param startDate      query string false "Data inicial para submissão da proposta. Formato aceito: YYYY-MM-DD"
// @Param endDate        query string false "Data final para submissão da proposta. Formato aceito: YYYY-MM-DD"
// @Param type           query string false "Tipo das matérias. Por padrão retorna todos os tipos. Valores permitidos: 'Proposição', 'Boletim'."
// @Param page           query int    false "Número da página. Por padrão é 1."
// @Param itemsPerPage   query int    false "Quantidade de matérias retornadas por página. Por padrão é 15."
// @Success 200 {array}  response.SwaggerNewsPagination "Requisição bem sucedida"
// @Failure 400 {object} response.SwaggerError          "Algum dado informado durante a requisição é inválido"
// @Failure 500 {object} response.SwaggerError          "Ocorreu um erro inesperado durante o processamento da requisição"
// @Router /news/trending [get]
func (instance News) GetTrending(context echo.Context) error {
	var newsFilter filters.NewsFilter

	newsFilter.Content = context.QueryParam("content")

	deputyIdParam := context.QueryParam("deputyId")
	if deputyIdParam != "" {
		deputyId, err := convertToUuid(deputyIdParam, "ID do deputado")
		if err != nil {
			log.Error(err)
			return context.JSON(http.StatusBadRequest, response.NewError(err.Error()))
		}
		newsFilter.DeputyId = &deputyId
	}

	partyIdParam := context.QueryParam("partyId")
	if partyIdParam != "" {
		partyId, err := convertToUuid(partyIdParam, "ID do partido")
		if err != nil {
			log.Error(err)
			return context.JSON(http.StatusBadRequest, response.NewError(err.Error()))
		}
		newsFilter.PartyId = &partyId
	}

	organizationIdParam := context.QueryParam("organizationId")
	if organizationIdParam != "" {
		organizationId, err := convertToUuid(organizationIdParam, "ID do organização")
		if err != nil {
			log.Error(err)
			return context.JSON(http.StatusBadRequest, response.NewError(err.Error()))
		}
		newsFilter.OrganizationId = &organizationId
	}

	startDateParam := context.QueryParam("startDate")
	if startDateParam != "" {
		startDate, err := convertToTime(startDateParam, "Data inicial")
		if err != nil {
			log.Error(err)
			return context.JSON(http.StatusBadRequest, response.NewError(err.Error()))
		}
		newsFilter.StartDate = &startDate
	}

	endDateParam := context.QueryParam("endDate")
	if endDateParam != "" {
		endDate, err := convertToTime(endDateParam, "Data final")
		if err != nil {
			log.Error(err)
			return context.JSON(http.StatusBadRequest, response.NewError(err.Error()))
		}
		newsFilter.EndDate = &endDate
	}

	typeParam := context.QueryParam("type")
	if typeParam != "" {
		if typeParam != "Proposição" && typeParam != "Boletim" {
			err := errors.New(fmt.Sprintf("Parâmetro inválido: type"))
			log.Errorf("Requisição mal formulada: %s (Valor: %s)", err, typeParam)
			return context.JSON(http.StatusBadRequest, response.NewError(err.Error()))
		} else if typeParam == "Boletim" && (newsFilter.DeputyId != nil || newsFilter.PartyId != nil || newsFilter.OrganizationId != nil) {
			err := errors.New(fmt.Sprintf("Parâmetro inválido: type igual 'Boletim' não permite o uso dos demais parâmetros query deputyId, partyId e organizationId."))
			log.Error("Parâmetro inválido: type igual 'Boletim' não permite o uso dos demais parâmetros query deputyId, partyId e organizationId.")
			return context.JSON(http.StatusBadRequest, response.NewError(err.Error()))
		}
		newsFilter.Type = typeParam
	}

	pageParam := context.QueryParam("page")
	if pageParam != "" {
		page, err := convertToInt(pageParam, "Página")
		if err != nil {
			log.Error(err)
			return context.JSON(http.StatusBadRequest, response.NewError(err.Error()))
		}
		newsFilter.PaginationFilter.Page = &page
	}

	itemsPerPageParam := context.QueryParam("itemsPerPage")
	if itemsPerPageParam != "" {
		perPage, err := convertToInt(itemsPerPageParam, "Itens por página")
		if err != nil {
			log.Error(err)
			return context.JSON(http.StatusBadRequest, response.NewError(err.Error()))
		}
		newsFilter.PaginationFilter.ItemsPerPage = &perPage
	}

	newsList, totalNumberOfNews, err := instance.service.GetTrending(newsFilter)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, response.NewError(err.Error()))
	}

	var newsResponse []response.News
	for _, news := range newsList {
		newsResponse = append(newsResponse, *response.NewNews(news))
	}

	requestResult := response.Pagination{
		Page:         newsFilter.PaginationFilter.GetPage(),
		ItensPerPage: newsFilter.PaginationFilter.GetItemsPerPage(),
		Total:        totalNumberOfNews,
		Data:         newsResponse,
	}

	return context.JSON(http.StatusOK, requestResult)
}
