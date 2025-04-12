package redis

import (
	"github.com/google/uuid"
)

type Session interface {
	CreateSession(userId uuid.UUID, sessionId uuid.UUID, accessToken string, refreshToken string) error
	SessionExists(userId uuid.UUID, sessionId uuid.UUID, accessToken string) (bool, error)
	RefreshTokenExists(userId uuid.UUID, sessionId uuid.UUID, accessToken string) (bool, error)
	DeleteSession(userId uuid.UUID, sessionId uuid.UUID) error
	DeleteSessionsByUserId(userId uuid.UUID) error
}
