package diconteiner

import "vnc-read-api/api/endpoints/handlers"

func GetPropositionHandler() *handlers.Proposition {
	return handlers.NewPropositionHandler(GetPropositionService())
}
