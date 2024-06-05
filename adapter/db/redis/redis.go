package redis_adapter

import "github.com/redis/go-redis/v9"

func Connect() *redis.Client {
	return redis.NewClient(&redis.Options{})
}
