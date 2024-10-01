package router

import (
	"github.com/labstack/echo/v4"
	"vnc-api/api/config/diconteiner"
)

func loadArticleRoutes(group *echo.Group) {
	newsHandler := diconteiner.GetArticleHandler()

	group = group.Group("/articles")

	group.GET("", newsHandler.GetArticles)
	group.GET("/trending", newsHandler.GetTrendingArticles)
	group.GET("/trending/type", newsHandler.GetTrendingArticlesByTypeId)
	group.GET("/view-later", newsHandler.GetArticlesToViewLater)
	group.GET("/:articleId/proposition", newsHandler.GetPropositionArticleById)
	group.GET("/:articleId/newsletter", newsHandler.GetNewsletterArticleById)
	group.PUT("/:articleId/rating", newsHandler.SaveArticleRating)
	group.PUT("/:articleId/view-later", newsHandler.SaveArticleToViewLater)
}
