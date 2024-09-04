package router

import (
	"github.com/labstack/echo/v4"
	"vnc-api/api/config/diconteiner"
)

func loadUserRoutes(group *echo.Group) {
	userHandler := diconteiner.GetUserHandler()

	group = group.Group("/user")

	group.PATCH("/resend-activation-email", userHandler.ResendActivationEmail)
	group.PATCH("/activate-account", userHandler.ActivateAccount)
}
