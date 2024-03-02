package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Hymiside/hezzl-api/pkg/custerrors"
	"github.com/Hymiside/hezzl-api/pkg/models"
	"github.com/redis/go-redis/v9"
)

type RepositoryRedis struct {
	db *redis.Client
}

func NewRepositoryRedis(db *redis.Client) *RepositoryRedis {
	return &RepositoryRedis{db: db}
}

func (r *RepositoryRedis) Set(ctx context.Context, b []byte) error {
	if err := r.db.Set(ctx, "goods", b, 1*time.Minute).Err(); err != nil {
		return fmt.Errorf("error to set goods: %v", err)
	}
	return nil
}

func (r *RepositoryRedis) Get(ctx context.Context, limit, offset int) (models.GoodsResponse, error) {
	b, goods := []byte{}, []models.Good{}
	if err := r.db.Get(ctx, "goods").Scan(&b); err != nil {
		if errors.Is(err, redis.Nil) {
			return models.GoodsResponse{}, custerrors.ErrNotFound
		}
		return models.GoodsResponse{}, fmt.Errorf("error to get goods: %v", err)
	}

	if err := json.Unmarshal(b, &goods); err != nil {
		return models.GoodsResponse{}, fmt.Errorf("error to unmarshal goods: %v", err)
	}

	var countRemoved = 0
	for i := 0; i < len(goods); i++ {
		if goods[i].Removed {
			countRemoved++
		}
	}

	goodsResponse := models.GoodsResponse{
		Goods: goods[offset:offset+limit],
		Meta: models.Meta{
			Limit:   limit,
			Ofset:   offset,
			Total:   len(goods),
			Removed: countRemoved,
		},
	}
	return goodsResponse, nil
}

func (r *RepositoryRedis) Delete(ctx context.Context) error {
	if err := r.db.Del(ctx, "goods").Err(); err != nil {
		return fmt.Errorf("error to delete goods: %v", err)
	}
	return nil
}