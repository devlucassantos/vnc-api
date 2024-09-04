package services

import (
	"github.com/devlucassantos/vnc-domains/src/domains/user"
)

type Email interface {
	SendUserAccountActivationEmail(userData user.User) error
}
