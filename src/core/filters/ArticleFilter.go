package filters

import (
	"github.com/google/uuid"
	"time"
)

type ArticleFilter struct {
	TypeId           *uuid.UUID
	Content          string
	DeputyId         *uuid.UUID
	PartyId          *uuid.UUID
	ExternalAuthorId *uuid.UUID
	StartDate        *time.Time
	EndDate          *time.Time
	PaginationFilter PaginationFilter
}
