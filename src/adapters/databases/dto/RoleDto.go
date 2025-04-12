package dto

import (
	"github.com/google/uuid"
)

type Role struct {
	Id   uuid.UUID `db:"role_id"`
	Code string    `db:"role_code"`
}
