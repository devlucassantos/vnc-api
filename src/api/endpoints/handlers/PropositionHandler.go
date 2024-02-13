package handlers

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"strings"
	"vnc-read-api/api/endpoints/dto/response"
	"vnc-read-api/core/interfaces/services"
)

type Proposition struct {
	propositionService services.Proposition
	newsletterService  services.Newsletter
}

func NewPropositionHandler(propositionService services.Proposition, newsletterService services.Newsletter) *Proposition {
	return &Proposition{
		propositionService: propositionService,
		newsletterService:  newsletterService,
	}
}

// GetPropositionById
// @ID 			GetPropositionById
// @Summary 	Busca dos detalhes da proposição pelo ID
// @Tags 		Proposições
// @Description Esta requisição é responsável por retornar os detalhes da proposição pelo ID.
// @Produce		json
// @Param propositionId path string true "ID da proposição"
// @Success 200 {array}  response.SwaggerProposition "Requisição bem sucedida"
// @Failure 400 {object} response.SwaggerError       "Algum dado informado durante a requisição é inválido"
// @Failure 500 {object} response.SwaggerError       "Ocorreu um erro inesperado durante o processamento da requisição"
// @Router /news/propositions/{propositionId} [get]
func (instance Proposition) GetPropositionById(context echo.Context) error {
	propositionId, err := uuid.Parse(context.Param("propositionId"))
	if err != nil {
		log.Error("Requisição mal formulada: Parâmetro inválido: ID da proposição - ", err.Error())
		return context.JSON(http.StatusBadRequest, response.NewError(fmt.Sprintf("Parâmetro inválido: ID da proposição")))
	}

	proposition, err := instance.propositionService.GetPropositionById(propositionId)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return context.JSON(http.StatusNotFound, response.NewError(fmt.Sprintf("Proposição não encontrada")))
		}
		return context.JSON(http.StatusInternalServerError, response.NewError(err.Error()))
	}

	propositionResponse := response.NewProposition(*proposition)

	newsletter, _ := instance.newsletterService.GetNewsletterByPropositionId(propositionId)
	if newsletter != nil {
		propositionResponse.Newsletter = response.NewNewsletter(*newsletter)
	}

	return context.JSON(http.StatusOK, propositionResponse)
}
