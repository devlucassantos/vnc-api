package handlers

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"strings"
	"vnc-read-api/api/endpoints/dto/filter"
	"vnc-read-api/api/endpoints/dto/response"
	"vnc-read-api/core/interfaces/services"
)

type Proposition struct {
	service services.Proposition
}

func NewPropositionHandler(service services.Proposition) *Proposition {
	return &Proposition{service: service}
}

// GetPropositions
// @ID 			GetPropositions
// @Summary 	Listar as proposições mais recentes
// @Tags 		Proposições
// @Description Esta requisição é responsável por retornar as proposições mais recentes disponíveis na plataforma Você na Câmara.
// @Produce		json
// @Param page         query int false "Número da página. Por padrão é 1."
// @Param itemsPerPage query int false "Quantidade de proposições retornadas por página. Por padrão é 25."
// @Success 200 {array}  response.SwaggerPropositionPagination "Requisição bem sucedida"
// @Failure 400 {object} response.SwaggerError                 "Algum dado informado durante a requisição é inválido"
// @Failure 500 {object} response.SwaggerError                 "Ocorreu um erro inesperado durante o processamento da requisição"
// @Router /proposition [get]
func (instance Proposition) GetPropositions(context echo.Context) error {
	var propositionFilter filter.PropositionFilter
	pageParam := context.QueryParam("page")
	if pageParam != "" {
		page, err := convertToInt(pageParam, "Página")
		if err != nil {
			log.Error(err)
			return context.JSON(http.StatusBadRequest, response.NewError(err.Error()))
		}
		propositionFilter.PaginationFilter.Page = &page
	}

	itemsPerPageParam := context.QueryParam("itemsPerPage")
	if itemsPerPageParam != "" {
		perPage, err := convertToInt(itemsPerPageParam, "Itens por página")
		if err != nil {
			log.Error(err)
			return context.JSON(http.StatusBadRequest, response.NewError(err.Error()))
		}
		propositionFilter.PaginationFilter.ItemsPerPage = &perPage
	}

	propositions, totalNumberOfPropositions, err := instance.service.GetPropositions(propositionFilter)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, response.NewError(err.Error()))
	}

	var propositionsResponse []response.Proposition
	for _, proposition := range propositions {
		propositionsResponse = append(propositionsResponse, response.NewProposition(proposition))
	}

	requestResult := response.Pagination{
		Page:         propositionFilter.PaginationFilter.GetPage(),
		ItensPerPage: propositionFilter.PaginationFilter.GetItemsPerPage(),
		Total:        totalNumberOfPropositions,
		Data:         propositionsResponse,
	}

	return context.JSON(http.StatusOK, requestResult)
}

// GetPropositionById
// @ID 			GetPropositionById
// @Summary 	Buscar os detalhes de uma proposição pelo seu ID
// @Tags 		Proposições
// @Description Esta requisição é responsável por retornar os detalhes de uma proposição pelo seu ID.
// @Produce		json
// @Param propositionId path string true "ID da proposição"
// @Success 200 {array}  response.SwaggerPropositionPagination "Requisição bem sucedida"
// @Failure 400 {object} response.SwaggerError                 "Algum dado informado durante a requisição é inválido"
// @Failure 500 {object} response.SwaggerError                 "Ocorreu um erro inesperado durante o processamento da requisição"
// @Router /proposition/{propositionId} [get]
func (instance Proposition) GetPropositionById(context echo.Context) error {
	propositionId, err := uuid.Parse(context.Param("propositionId"))
	if err != nil {
		log.Error("Requisição mal formulada: Parâmetro inválido: ID da proposição - ", err.Error())
		return context.JSON(http.StatusBadRequest, response.NewError(fmt.Sprintf("Parâmetro inválido: ID da proposição")))
	}

	proposition, err := instance.service.GetPropositionById(propositionId)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return context.JSON(http.StatusNotFound, response.NewError(fmt.Sprintf("Proposição não encontrada")))
		}
		return context.JSON(http.StatusInternalServerError, response.NewError(err.Error()))
	}

	return context.JSON(http.StatusOK, response.NewProposition(*proposition))
}
