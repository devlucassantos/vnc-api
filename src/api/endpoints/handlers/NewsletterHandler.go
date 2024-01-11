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

type Newsletter struct {
	service services.Newsletter
}

func NewNewsletterHandler(service services.Newsletter) *Newsletter {
	return &Newsletter{service: service}
}

// GetNewsletterById
// @ID 			GetNewsletterById
// @Summary 	Busca dos detalhes do boletim pelo ID
// @Tags 		Boletins
// @Description Esta requisição é responsável por retornar os detalhes do boletim pelo ID.
// @Produce		json
// @Param newsletterId path string true "ID do boletim"
// @Success 200 {array}  response.SwaggerNewsletter "Requisição bem sucedida"
// @Failure 400 {object} response.SwaggerError      "Algum dado informado durante a requisição é inválido"
// @Failure 500 {object} response.SwaggerError      "Ocorreu um erro inesperado durante o processamento da requisição"
// @Router /news/newsletters/{newsletterId} [get]
func (instance Newsletter) GetNewsletterById(context echo.Context) error {
	newsletterId, err := uuid.Parse(context.Param("newsletterId"))
	if err != nil {
		log.Error("Requisição mal formulada: Parâmetro inválido: ID da proposição - ", err.Error())
		return context.JSON(http.StatusBadRequest, response.NewError(fmt.Sprintf("Parâmetro inválido: ID da proposição")))
	}

	newsletter, err := instance.service.GetNewsletterById(newsletterId)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return context.JSON(http.StatusNotFound, response.NewError(fmt.Sprintf("Proposição não encontrada")))
		}
		return context.JSON(http.StatusInternalServerError, response.NewError(err.Error()))
	}

	return context.JSON(http.StatusOK, response.NewNewsletter(*newsletter))
}
