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
		log.Error("Error creating a connection to the Redis database: ", err.Error())
		return err
	}
	defer instance.connectionManager.closeConnection(redisConnection)

	err = redisConnection.Set(fmt.Sprintf("access:%s:%s", userId, sessionId), accessToken, user.AccessTokenTimeout).Err()
	if err != nil {
		log.Errorf("Error registering access token for user %s (Session: %s): %s", userId, sessionId, err.Error())
		return err
	}

	err = redisConnection.Set(fmt.Sprintf("refresh:%s:%s", userId, sessionId), refreshToken, user.RefreshTokenTimeout).Err()
	if err != nil {
		log.Errorf("Error registering refresh token for user %s (Session: %s): %s", userId, sessionId, err.Error())
		return err
	}

	return nil
}

func (instance Session) SessionExists(userId uuid.UUID, sessionId uuid.UUID, handledToken string) (bool, error) {
	redisConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("Error creating a connection to the Redis database: ", err.Error())
		return false, err
	}
	defer instance.connectionManager.closeConnection(redisConnection)

	accessKey := fmt.Sprintf("access:%s:%s", userId, sessionId)
	keyExists, err := redisConnection.Exists(accessKey).Result()
	if err != nil {
		log.Errorf("Error checking the existence of the access key for user %s (Session: %s): %s",
			userId, sessionId, err.Error())
		return false, err
	}

	if keyExists <= 0 {
		return false, nil
	}

	storedValue, err := redisConnection.Get(accessKey).Result()
	if err != nil {
		log.Errorf("Error fetching access key value for user %s (Session: %s): %s", userId, sessionId, err.Error())
		return false, err
	}

	return storedValue == handledToken, nil
}

func (instance Session) RefreshTokenExists(userId uuid.UUID, sessionId uuid.UUID, handledToken string) (bool, error) {
	redisConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("Error creating a connection to the Redis database: ", err.Error())
		return false, err
	}
	defer instance.connectionManager.closeConnection(redisConnection)

	refreshKey := fmt.Sprintf("refresh:%s:%s", userId, sessionId)
	keyExists, err := redisConnection.Exists(refreshKey).Result()
	if err != nil {
		log.Errorf("Error checking the existence of the refresh key for user %s (Session: %s): %s",
			userId, sessionId, err.Error())
		return false, err
	}

	if keyExists <= 0 {
		return false, nil
	}

	storedValue, err := redisConnection.Get(refreshKey).Result()
	if err != nil {
		log.Errorf("Error fetching refresh key value for user %s (Session: %s): %s", userId, sessionId, err.Error())
		return false, err
	}

	return storedValue == handledToken, nil
}

func (instance Session) DeleteSession(userId uuid.UUID, sessionId uuid.UUID) error {
	redisConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("Error creating a connection to the Redis database: ", err.Error())
		return err
	}
	defer instance.connectionManager.closeConnection(redisConnection)

	err = redisConnection.Del(fmt.Sprintf("access:%s:%s", userId, sessionId),
		fmt.Sprintf("refresh:%s:%s", userId, sessionId)).Err()
	if err != nil {
		log.Errorf("Error deleting the session %s for user %s: %s", sessionId, userId, err.Error())
		return err
	}

	return nil
}

func (instance Session) DeleteSessionsByUserId(userId uuid.UUID) error {
	redisConnection, err := instance.connectionManager.createConnection()
	if err != nil {
		log.Error("Error creating a connection to the Redis database: ", err.Error())
		return err
	}
	defer instance.connectionManager.closeConnection(redisConnection)

	iterator := redisConnection.Scan(0, fmt.Sprintf("*:%s:*", userId), 0).Iterator()

	for iterator.Next() {
		err = redisConnection.Del(iterator.Val()).Err()
		if err != nil {
			log.Errorf("Error deleting the sessions for user %s: %s", userId, err.Error())
			return err
		}
	}

	err = iterator.Err()
	if err != nil {
		log.Errorf("Error during the iteration responsible for deleting the sessions of user %s: %s", userId, err.Error())
		return err
	}

	return nil
}
