package dicontainer

import (
	"vnc-api/adapters/api/endpoints/handlers"
)

func GetAuthenticationHandler() *handlers.Authentication {
	return handlers.NewAuthenticationHandler(GetAuthenticationService())
}

func GetUserHandler() *handlers.User {
	return handlers.NewUserHandler(GetUserService())
}

func GetResourcesHandler() *handlers.Resources {
	return handlers.NewResourcesHandler(GetResourcesService())
}

func GetArticleHandler() *handlers.Article {
	return handlers.NewArticleHandler(GetArticleService(), GetResourcesService(), GetPropositionService(),
		GetVotingService(), GetEventService(), GetNewsletterService())
}
