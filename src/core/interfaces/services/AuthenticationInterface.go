package services

import (
	"github.com/devlucassantos/vnc-domains/src/domains/user"
	"github.com/google/uuid"
)

type Authentication interface {
	SignUp(signUpData user.User) (*user.User, error)
	SignIn(signInData user.User) (*user.User, error)
	SignOut(userId uuid.UUID, sessionId uuid.UUID) error
	RefreshTokens(userId uuid.UUID, sessionId uuid.UUID, refreshToken string) (*user.User, error)
	SessionExists(userId uuid.UUID, sessionId uuid.UUID, token string) (bool, error)
}
