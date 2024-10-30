package redis

import (
	"crypto/tls"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/labstack/gommon/log"
	"os"
	"strings"
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
	redisOptions, err := getRedisConnectionOptions()
	if err != nil {
		log.Error("Error getting Redis database connection address: ", err.Error())
		return nil, err
	}

	connection := redis.NewClient(redisOptions)

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

func getRedisConnectionOptions() (*redis.Options, error) {
	redisUrl := os.Getenv("REDIS_URL")
	if len(redisUrl) > 0 {
		redisOptions, err := redis.ParseURL(redisUrl)
		if err != nil {
			log.Error("Error parsing Redis URL: ", err.Error())
			return nil, err
		}

		if strings.HasPrefix(redisUrl, "rediss") {
			redisOptions.TLSConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}

		return redisOptions, nil
	}

	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")

	redisOptions := &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       0,
	}

	return redisOptions, nil
}
