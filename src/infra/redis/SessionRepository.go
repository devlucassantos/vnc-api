package redis

import (
	"fmt"
	"github.com/devlucassantos/vnc-domains/src/domains/user"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
)

type Session struct {
	connectionManager connectionManagerInterface
}

func NewSessionRepository(connectionManager connectionManagerInterface) *Session {
	return &Session{
		connectionManager: connectionManager,
	}
}

func (instance Session) CreateSession(userId uuid.UUID, sessionId uuid.UUID, accessToken string, refreshToken string) error {
	redisConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("Erro ao tentar se conectar com o Redis: ", err.Error())
		return err
	}
	defer instance.connectionManager.closeConnection(redisConnection)

	err = redisConnection.Set(fmt.Sprintf("access:%s:%s", userId, sessionId), accessToken, user.AccessTokenTimeout).Err()
	if err != nil {
		log.Errorf("Erro ao registrar token de acesso do usuário %s: ", userId, err.Error())
		return err
	}

	err = redisConnection.Set(fmt.Sprintf("refresh:%s:%s", userId, sessionId), refreshToken, user.RefreshTokenTimeout).Err()
	if err != nil {
		log.Errorf("Erro ao registrar token de atualização do usuário %s: ", userId, err.Error())
		return err
	}

	return nil
}

func (instance Session) SessionExists(userId uuid.UUID, sessionId uuid.UUID, handledToken string) (bool, error) {
	redisConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("Erro ao tentar se conectar com o Redis: ", err.Error())
		return false, err
	}
	defer instance.connectionManager.closeConnection(redisConnection)

	accessKey := fmt.Sprintf("access:%s:%s", userId, sessionId)
	keyExists, err := redisConnection.Exists(accessKey).Result()
	if err != nil {
		log.Errorf("Erro ao verificar existência da chave de acesso do usuário %s: %s", userId, err.Error())
		return false, err
	}

	if keyExists <= 0 {
		return false, nil
	}

	storedValue, err := redisConnection.Get(accessKey).Result()
	if err != nil {
		log.Error("Erro ao buscar o valor da chave de acesso: ", err.Error())
		return false, err
	}

	return storedValue == handledToken, nil
}

func (instance Session) RefreshTokenExists(userId uuid.UUID, sessionId uuid.UUID, handledToken string) (bool, error) {
	redisConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("Erro ao tentar se conectar com o Redis: ", err.Error())
		return false, err
	}
	defer instance.connectionManager.closeConnection(redisConnection)

	refreshKey := fmt.Sprintf("refresh:%s:%s", userId, sessionId)
	keyExists, err := redisConnection.Exists(refreshKey).Result()
	if err != nil {
		log.Errorf("Erro ao verificar existência da chave de atualização do usuário %s: %s", userId, err.Error())
		return false, err
	}

	if keyExists <= 0 {
		return false, nil
	}

	storedValue, err := redisConnection.Get(refreshKey).Result()
	if err != nil {
		log.Error("Erro ao buscar o valor da chave de atualização: ", err.Error())
		return false, err
	}

	return storedValue == handledToken, nil
}

func (instance Session) DeleteSession(userId uuid.UUID, sessionId uuid.UUID) error {
	redisConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("Erro ao tentar se conectar com o Redis: ", err.Error())
		return err
	}
	defer instance.connectionManager.closeConnection(redisConnection)

	result := redisConnection.Del(fmt.Sprintf("access:%s:%s", userId, sessionId),
		fmt.Sprintf("refresh:%s:%s", userId, sessionId))
	if result.Err() != nil {
		log.Errorf("Erro ao deletar sessão do usuário %s: %s", userId, err)
		return err
	}

	return nil
}
