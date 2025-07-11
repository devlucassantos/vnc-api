package router

import (
	"github.com/labstack/echo/v4"
	"vnc-api/config/dicontainer"
)

func loadUserRoutes(group *echo.Group) {
	userHandler := dicontainer.GetUserHandler()

	group = group.Group("/user")

	group.PATCH("/resend-activation-email", userHandler.ResendActivationEmail)
	group.PATCH("/activate-account", userHandler.ActivateAccount)
}
