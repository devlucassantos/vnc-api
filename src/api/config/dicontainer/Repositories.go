package dicontainer

import (
	"vnc-api/core/interfaces/repositories"
	"vnc-api/infra/postgres"
	"vnc-api/infra/redis"
)

func GetRedisDatabaseManager() *redis.ConnectionManager {
	return redis.NewRedisConnectionManager()
}

func GetPostgresDatabaseManager() *postgres.ConnectionManager {
	return postgres.NewPostgresConnectionManager()
}

func GetSessionRedisRepository() repositories.Session {
	return redis.NewSessionRepository(GetRedisDatabaseManager())
}

func GetUserPostgresRepository() repositories.User {
	return postgres.NewUserRepository(GetPostgresDatabaseManager())
}

func GetResourcesPostgresRepository() repositories.Resources {
	return postgres.NewResourcesRepository(GetPostgresDatabaseManager())
}

func GetArticlePostgresRepository() repositories.Article {
	return postgres.NewArticleRepository(GetPostgresDatabaseManager())
}

func GetPropositionPostgresRepository() repositories.Proposition {
	return postgres.NewPropositionRepository(GetPostgresDatabaseManager())
}

func GetVotingPostgresRepository() repositories.Voting {
	return postgres.NewVotingRepository(GetPostgresDatabaseManager())
}

func GetEventPostgresRepository() repositories.Event {
	return postgres.NewEventRepository(GetPostgresDatabaseManager())
}

func GetNewsletterPostgresRepository() repositories.Newsletter {
	return postgres.NewNewsletterRepository(GetPostgresDatabaseManager())
}
