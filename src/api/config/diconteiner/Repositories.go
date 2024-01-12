package diconteiner

import (
	"vnc-read-api/core/interfaces/repositories"
	"vnc-read-api/infra/postgres"
)

func GetPostgresDatabaseManager() *postgres.ConnectionManager {
	return postgres.NewPostgresConnectionManager()
}

func GetResourcesPostgresRepository() repositories.Resources {
	return postgres.NewResourcesRepository(GetPostgresDatabaseManager())
}

func GetPropositionPostgresRepository() repositories.Proposition {
	return postgres.NewPropositionRepository(GetPostgresDatabaseManager())
}

func GetNewsletterPostgresRepository() repositories.Newsletter {
	return postgres.NewNewsletterRepository(GetPostgresDatabaseManager())
}

func GetNewsPostgresRepository() repositories.News {
	return postgres.NewNewsRepository(GetPostgresDatabaseManager())
}
