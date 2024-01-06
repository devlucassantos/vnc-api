package diconteiner

import (
	interfaces "vnc-read-api/core/interfaces/services"
	"vnc-read-api/core/services"
)

func GetPropositionService() interfaces.Proposition {
	return services.NewPropositionService(GetPropositionPostgresRepository())
}

func GetNewsletterService() interfaces.Newsletter {
	return services.NewNewsletterService(GetNewsletterPostgresRepository())
}

func GetNewsService() interfaces.News {
	return services.NewNewsService(GetNewsPostgresRepository())
}
