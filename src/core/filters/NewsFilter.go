package filters

import (
	"github.com/google/uuid"
	"time"
)

type NewsFilter struct {
	Content          string
	DeputyId         *uuid.UUID
	PartyId          *uuid.UUID
	OrganizationId   *uuid.UUID
	StartDate        *time.Time
	EndDate          *time.Time
	Type             string
	PaginationFilter PaginationFilter
}
