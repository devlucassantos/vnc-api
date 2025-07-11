package dicontainer

import (
	"vnc-api/adapters/databases/postgres"
	interfaces "vnc-api/core/interfaces/postgres"
)

func GetPostgresDatabaseManager() *postgres.ConnectionManager {
	return postgres.NewPostgresConnectionManager()
}

func GetUserPostgresRepository() interfaces.User {
	return postgres.NewUserRepository(GetPostgresDatabaseManager())
}

func GetResourcesPostgresRepository() interfaces.Resources {
	return postgres.NewResourcesRepository(GetPostgresDatabaseManager())
}

func GetArticlePostgresRepository() interfaces.Article {
	return postgres.NewArticleRepository(GetPostgresDatabaseManager())
}

func GetPropositionPostgresRepository() interfaces.Proposition {
	return postgres.NewPropositionRepository(GetPostgresDatabaseManager())
}

func GetVotingPostgresRepository() interfaces.Voting {
	return postgres.NewVotingRepository(GetPostgresDatabaseManager())
}

func GetEventPostgresRepository() interfaces.Event {
	return postgres.NewEventRepository(GetPostgresDatabaseManager())
}

func GetNewsletterPostgresRepository() interfaces.Newsletter {
	return postgres.NewNewsletterRepository(GetPostgresDatabaseManager())
}
