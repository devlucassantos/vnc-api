package routes

import (
	"github.com/labstack/echo/v4"
	"vnc-read-api/api/config/diconteiner"
)

func loadNewsRoutes(group *echo.Group) {
	newsHandler := diconteiner.GetNewsHandler()

	group = group.Group("/news")
	loadPropositionRoutes(group)
	loadNewsletterRoutes(group)

	group.GET("", newsHandler.GetNews)
	group.GET("/trending", newsHandler.GetTrending)
}
