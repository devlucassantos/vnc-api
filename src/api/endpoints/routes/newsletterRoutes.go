package routes

import (
	"github.com/labstack/echo/v4"
	"vnc-read-api/api/config/diconteiner"
)

func loadNewsletterRoutes(group *echo.Group) {
	newsletterHandler := diconteiner.GetNewsletterHandler()

	group.GET("/newsletters/:newsletterId", newsletterHandler.GetNewsletterById)
}
