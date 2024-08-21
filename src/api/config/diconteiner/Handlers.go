package diconteiner

import "vnc-api/api/endpoints/handlers"

func GetAuthenticationHandler() *handlers.Authentication {
	return handlers.NewAuthenticationHandler(GetAuthenticationService())
}

func GetResourcesHandler() *handlers.Resources {
	return handlers.NewResourcesHandler(GetResourcesService())
}

func GetArticleHandler() *handlers.Article {
	return handlers.NewArticleHandler(GetArticleService(), GetResourcesService(), GetPropositionService(), GetNewsletterService())
}
