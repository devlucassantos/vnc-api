package routes

import (
	"github.com/labstack/echo/v4"
	"vnc-read-api/api/config/diconteiner"
)

func loadResourcesRoutes(group *echo.Group) {
	resourcesHandler := diconteiner.GetResourcesHandler()

	group = group.Group("/resources")

	group.GET("", resourcesHandler.GetResources)
}
