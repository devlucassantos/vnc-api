package diconteiner

import "vnc-read-api/api/endpoints/handlers"

func GetPropositionHandler() *handlers.Proposition {
	return handlers.NewPropositionHandler(GetPropositionService())
}

func GetNewsletterHandler() *handlers.Newsletter {
	return handlers.NewNewsletterHandler(GetNewsletterService())
}

func GetNewsHandler() *handlers.News {
	return handlers.NewNewsHandler(GetNewsService())
}
