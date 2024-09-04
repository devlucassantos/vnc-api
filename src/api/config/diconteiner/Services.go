package diconteiner

import (
	interfaces "vnc-api/core/interfaces/services"
	"vnc-api/core/services"
)

func GetAuthenticationService() interfaces.Authentication {
	return services.NewAuthenticationService(GetUserPostgresRepository(), GetSessionRedisRepository(), GetEmailService())
}

func GetUserService() interfaces.User {
	return services.NewUserService(GetUserPostgresRepository(), GetSessionRedisRepository(), GetEmailService())
}

func GetResourcesService() interfaces.Resources {
	return services.NewResourcesService(GetResourcesPostgresRepository())
}

func GetArticleService() interfaces.Article {
	return services.NewArticleService(GetArticlePostgresRepository())
}

func GetPropositionService() interfaces.Proposition {
	return services.NewPropositionService(GetPropositionPostgresRepository())
}

func GetNewsletterService() interfaces.Newsletter {
	return services.NewNewsletterService(GetNewsletterPostgresRepository())
}

func GetEmailService() interfaces.Email {
	return services.NewEmailService()
}
