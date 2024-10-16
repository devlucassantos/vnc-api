package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/labstack/gommon/log"
	"os"
)

type connectionManagerInterface interface {
	createConnection() (*redis.Client, error)
	closeConnection(connection *redis.Client)
}

type ConnectionManager struct{}

func NewRedisConnectionManager() *ConnectionManager {
	return &ConnectionManager{}
}

func (ConnectionManager) createConnection() (*redis.Client, error) {
	connection := redis.NewClient(getRedisConnectionOptions())

	pingResult := connection.Ping()
	if pingResult.Err() != nil {
		log.Error("Error checking the creation of the connection to the Redis database: ", pingResult.Err().Error())
		return nil, pingResult.Err()
	}

	return connection, nil
}

func (ConnectionManager) closeConnection(connection *redis.Client) {
	err := connection.Close()
	if err != nil {
		log.Error("Error closing the connection to the Redis database: ", err.Error())
	}
}

func getRedisConnectionOptions() *redis.Options {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")

	redisOptions := &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       0,
	}

	return redisOptions
}
