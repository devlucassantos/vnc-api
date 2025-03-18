package filters

import (
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"time"
)

type Article struct {
	TypeId         *uuid.UUID
	SpecificTypeId *uuid.UUID
	Content        string
	StartDate      *time.Time
	EndDate        *time.Time
	Proposition
	Voting
	Event
	Pagination
}

func (instance Article) HasConflict() error {
	isPropositionFilterActive := !instance.Proposition.IsZero()
	isVotingFilterActive := !instance.Voting.IsZero()
	isEventFilterActive := !instance.Event.IsZero()

	var errorMessage string
	if isPropositionFilterActive && isVotingFilterActive && isEventFilterActive {
		errorMessage = "Invalid parameters: Conflicting parameters for propositions, votes and events"
	} else if isPropositionFilterActive && isVotingFilterActive {
		errorMessage = "Invalid parameters: Conflicting parameters for propositions and votes"
	} else if isPropositionFilterActive && isEventFilterActive {
		errorMessage = "Invalid parameters: Conflicting parameters for propositions and events"
	} else if isVotingFilterActive && isEventFilterActive {
		errorMessage = "Invalid parameters: Conflicting parameters for votes and events"
	}

	if errorMessage != "" {
		log.Warn("Conflict identified when checking filters: ", errorMessage)
		return errors.New(errorMessage)
	}

	return nil
}
