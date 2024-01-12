package diconteiner

import "vnc-read-api/api/endpoints/handlers"

func GetResourcesHandler() *handlers.Resources {
	return handlers.NewResourcesHandler(GetResourcesService())
}

func GetPropositionHandler() *handlers.Proposition {
	return handlers.NewPropositionHandler(GetPropositionService(), GetNewsletterService())
}

func GetNewsletterHandler() *handlers.Newsletter {
	return handlers.NewNewsletterHandler(GetNewsletterService())
}

func GetNewsHandler() *handlers.News {
	return handlers.NewNewsHandler(GetNewsService())
}
