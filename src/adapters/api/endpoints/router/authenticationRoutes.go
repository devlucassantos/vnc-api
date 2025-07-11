package router

import (
	"github.com/labstack/echo/v4"
	"vnc-api/config/dicontainer"
)

func loadAuthenticationRoutes(group *echo.Group) {
	authenticationHandler := dicontainer.GetAuthenticationHandler()

	group = group.Group("/auth")

	group.POST("/sign-up", authenticationHandler.SignUp)
	group.POST("/sign-in", authenticationHandler.SignIn)
	group.POST("/sign-out", authenticationHandler.SignOut)
	group.POST("/refresh", authenticationHandler.Refresh)
}
