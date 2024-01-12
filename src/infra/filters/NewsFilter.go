package filters

import (
	"database/sql"
	"github.com/google/uuid"
)

type NewsFilter struct {
	Content        sql.NullString
	DeputyId       uuid.NullUUID
	PartyId        uuid.NullUUID
	OrganizationId uuid.NullUUID
	StartDate      sql.NullTime
	EndDate        sql.NullTime
	Type           sql.NullString
	Offset         int
	Limit          int
}
