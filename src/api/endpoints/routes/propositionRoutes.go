package routes

import (
	"github.com/labstack/echo/v4"
	"vnc-read-api/api/config/diconteiner"
)

func loadPropositionRoutes(group *echo.Group) {
	propositionHandler := diconteiner.GetPropositionHandler()

	group.GET("/proposition", propositionHandler.GetPropositions)
	group.GET("/proposition/:propositionId", propositionHandler.GetPropositionById)
}
