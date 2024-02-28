package redis

import (
	"context"
	"fmt"

	"github.com/Hymiside/hezzl-api/pkg/models"
	"github.com/redis/go-redis/v9"
)

func NewRedisDB(ctx context.Context, cfg models.ConfigRedis) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: "",
		DB:       0,
	})
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("connection test error: %w", err)
	}
	return rdb, nil

}