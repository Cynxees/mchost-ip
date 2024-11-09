package queue

import (
	"github.com/redis/go-redis/v9"
)

func NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Connect to Redis inside Docker
		Password: "",               // No password set by default
		DB:       0,                // Use default DB
	})
}
