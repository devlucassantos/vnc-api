package dto

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	Id        uuid.UUID `db:"user_id"`
	FirstName string    `db:"user_first_name"`
	LastName  string    `db:"user_last_name"`
	Email     string    `db:"user_email"`
	Password  string    `db:"user_password"`
	CreatedAt time.Time `db:"user_created_at"`
	UpdatedAt time.Time `db:"user_updated_at"`
}
