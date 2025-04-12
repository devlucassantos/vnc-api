package router

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
)

func loadDocumentationRoutes(group *echo.Group) {
	group.GET("/documentation/*", echoSwagger.WrapHandler)
	group.GET("/documentation", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/api/documentation/index.html")
	})
}
