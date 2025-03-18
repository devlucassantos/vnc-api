package router

import (
	"github.com/labstack/echo/v4"
	"vnc-api/api/config/dicontainer"
)

func loadResourcesRoutes(group *echo.Group) {
	resourcesHandler := dicontainer.GetResourcesHandler()

	group = group.Group("/resources")

	group.GET("", resourcesHandler.GetResources)
}
