package dto

import (
	"github.com/google/uuid"
	"time"
)

type Role struct {
	Id        uuid.UUID `db:"role_id"`
	Code      string    `db:"role_code"`
	CreatedAt time.Time `db:"role_created_at"`
	UpdatedAt time.Time `db:"role_updated_at"`
}
