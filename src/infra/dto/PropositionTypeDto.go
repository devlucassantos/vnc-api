package dto

import (
	"github.com/google/uuid"
	"time"
)

type PropositionType struct {
	Id          uuid.UUID `db:"proposition_type_id"`
	Description string    `db:"proposition_type_description"`
	Color       string    `db:"proposition_type_color"`
	SortOrder   int       `db:"proposition_type_sort_order"`
	CreatedAt   time.Time `db:"proposition_type_created_at"`
	UpdatedAt   time.Time `db:"proposition_type_updated_at"`
}
