package dicontainer

import (
	"vnc-api/adapters/databases/redis"
	interfaces "vnc-api/core/interfaces/redis"
)

func GetRedisDatabaseManager() *redis.ConnectionManager {
	return redis.NewRedisConnectionManager()
}

func GetSessionRedisRepository() interfaces.Session {
	return redis.NewSessionRepository(GetRedisDatabaseManager())
}
