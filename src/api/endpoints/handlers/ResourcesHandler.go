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
// @Summary     List all resources
// @Tags        Resources
// @Description This request is responsible for listing all the platform resources.
// @Produce     json
// @Success 200 {array}  swagger.Resources "Successful request"
// @Failure 401 {object} swagger.HttpError "Unauthorized access"
// @Failure 500 {object} swagger.HttpError "An unexpected error occurred while processing the request"
// @Failure 503 {object} swagger.HttpError "Some of the services/resources are temporarily unavailable"
// @Router /resources [GET]
func (instance Resources) GetResources(context echo.Context) error {
	articleTypes, propositionTypes, parties, deputies, externalAuthors, legislativeBodies, eventTypes, eventSituations,
		err := instance.service.GetResources()
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			log.Error("Database unavailable: ", err.Error())
			return context.JSON(http.StatusServiceUnavailable, response.NewServiceUnavailableError())
		}

		log.Error("Error retrieving resource data: ", err.Error())
		return context.JSON(http.StatusInternalServerError, response.NewInternalServerError())
	}

	return context.JSON(http.StatusOK, response.NewResources(articleTypes, propositionTypes, parties, deputies,
		externalAuthors, legislativeBodies, eventTypes, eventSituations))
}
