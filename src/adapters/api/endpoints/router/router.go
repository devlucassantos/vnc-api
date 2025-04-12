package router

import (
	"github.com/labstack/echo/v4"
	"os"
)

type Router interface {
	Load(*echo.Group)
}

type router struct {
}

func New() Router {
	return &router{}
}

func (instance *router) Load(group *echo.Group) {
	if os.Getenv("SERVER_MODE") != "production" {
		loadDocumentationRoutes(group)
	}

	v1Group := group.Group("/v1")

	loadAuthenticationRoutes(v1Group)
	loadUserRoutes(v1Group)
	loadResourcesRoutes(v1Group)
	loadArticleRoutes(v1Group)
}
