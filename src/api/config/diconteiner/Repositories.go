package diconteiner

import (
	"vnc-read-api/core/interfaces/repositories"
	"vnc-read-api/infra/postgres"
)

func GetPostgresDatabaseManager() *postgres.ConnectionManager {
	return postgres.NewPostgresConnectionManager()
}

func GetPropositionPostgresRepository() repositories.Proposition {
	return postgres.NewPropositionRepository(GetPostgresDatabaseManager())
}
