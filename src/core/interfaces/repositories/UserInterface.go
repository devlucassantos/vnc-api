package repositories

import (
	"github.com/devlucassantos/vnc-domains/src/domains/user"
	"github.com/google/uuid"
)

type User interface {
	CreateUser(userData user.User) (*user.User, error)
	UpdateUser(userData user.User) (*user.User, error)
	GetUserById(id uuid.UUID) (*user.User, error)
	GetUserByEmail(email string) (*user.User, error)
}
