package diconteiner

import (
	interfaces "vnc-read-api/core/interfaces/services"
	"vnc-read-api/core/services"
)

func GetPropositionService() interfaces.Proposition {
	return services.NewPropositionService(GetPropositionPostgresRepository())
}
