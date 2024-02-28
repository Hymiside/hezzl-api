package redis

import "github.com/redis/go-redis/v9"

type RepositoryRedis struct {
	db *redis.Client
}

func NewRepositoryRedis(db *redis.Client) *RepositoryRedis {
	return &RepositoryRedis{db: db}
}