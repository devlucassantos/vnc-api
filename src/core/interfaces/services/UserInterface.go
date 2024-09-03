package services

import (
	"github.com/devlucassantos/vnc-domains/src/domains/user"
	"github.com/google/uuid"
)

type User interface {
	ResendUserAccountActivationEmail(userId uuid.UUID) error
	ActivateUserAccount(userAccountActivationData user.User) (*user.User, error)
}
