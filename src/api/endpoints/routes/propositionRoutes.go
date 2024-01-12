package routes

import (
	"github.com/labstack/echo/v4"
	"vnc-read-api/api/config/diconteiner"
)

func loadPropositionRoutes(group *echo.Group) {
	propositionHandler := diconteiner.GetPropositionHandler()

	group.GET("/propositions/:propositionId", propositionHandler.GetPropositionById)
}
