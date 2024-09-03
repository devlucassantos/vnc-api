package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"strings"
	"vnc-api/api/endpoints/dto/response"
	"vnc-api/core/interfaces/services"
)

type Resources struct {
	service services.Resources
}

func NewResourcesHandler(service services.Resources) *Resources {
	return &Resources{service: service}
}

// GetResources
// @ID          GetResources
// @Summary     Listar todos os recursos
// @Tags        Recursos
// @Description Esta requisição é responsável por listar todos os recursos da plataforma.
// @Produce     json
// @Success 200 {array}  response.SwaggerResources "Requisição realizada com sucesso."
// @Failure 401 {object} response.SwaggerHttpError "Acesso não autorizado."
// @Failure 500 {object} response.SwaggerHttpError "Ocorreu um erro inesperado durante o processamento da requisição."
// @Failure 503 {object} response.SwaggerHttpError "Algum dos serviços/recursos está temporariamente indisponível."
// @Router /resources [GET]
func (instance Resources) GetResources(context echo.Context) error {
	articleTypes, parties, deputies, externalAuthors, err := instance.service.GetResources()
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			log.Error("Banco de dados indisponível: ", err.Error())
			return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
		}

		log.Error("Erro ao buscar dados dos recursos no banco de dados: ", err.Error())
		return context.JSON(http.StatusInternalServerError, response.NewInternalServerError())
	}

	return context.JSON(http.StatusOK, response.NewResources(articleTypes, parties, deputies, externalAuthors))
}
