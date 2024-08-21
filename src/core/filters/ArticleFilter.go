package filters

import (
	"github.com/google/uuid"
	"time"
)

type ArticleFilter struct {
	Content          string
	DeputyId         *uuid.UUID
	PartyId          *uuid.UUID
	ExternalAuthorId *uuid.UUID
	StartDate        *time.Time
	EndDate          *time.Time
	Type             string
	PaginationFilter PaginationFilter
}
