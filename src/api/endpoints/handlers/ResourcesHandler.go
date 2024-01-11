package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"vnc-read-api/api/endpoints/dto/response"
	"vnc-read-api/core/interfaces/services"
)

type Resources struct {
	service services.Resources
}

func NewResourcesHandler(service services.Resources) *Resources {
	return &Resources{service: service}
}

// GetResources
// @ID 			GetResources
// @Summary 	Busca de todos os recursos
// @Tags 		Recursos
// @Description Esta requisição é responsável por retornar todos os recursos da plataforma.
// @Produce		json
// @Success 200 {array}  response.SwaggerResources "Requisição bem sucedida"
// @Failure 500 {object} response.SwaggerError     "Ocorreu um erro inesperado durante o processamento da requisição"
// @Router /resources [get]
func (instance Resources) GetResources(context echo.Context) error {
	parties, deputies, organizations, err := instance.service.GetResources()
	if err != nil {
		return context.JSON(http.StatusInternalServerError, response.NewError(err.Error()))
	}

	return context.JSON(http.StatusOK, response.NewResources(parties, deputies, organizations))
}
