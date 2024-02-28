package service

import (
	"context"
	"fmt"

	"github.com/Hymiside/hezzl-api/pkg/models"
)

type repositoryClickhouse interface{}

type repositoryPostgres interface{
	Goods(ctx context.Context, limit, offset int) (models.GoodsResponse, error)
	CreateGood(ctx context.Context, projectID int, name string) (models.Good, error)
	UpdateGood(ctx context.Context, good models.Good, goodID, projectID int) (models.Good, error)
}

type repositoryRedis interface{}

type queueNats interface{}

type Service struct {
	repoClickhouse repositoryClickhouse
	repoPostgres   repositoryPostgres
	repoRedis      repositoryRedis
	queueNats      queueNats
}

func NewService(
	repoClickhouse repositoryClickhouse,
	repoPostgres repositoryPostgres,
	repoRedis repositoryRedis,
	queueNats queueNats,
) *Service {
	return &Service{
		repoClickhouse: repoClickhouse,
		repoPostgres:   repoPostgres,
		repoRedis:      repoRedis,
		queueNats:      queueNats,
	}
}

func (s *Service) Goods(ctx context.Context, limit, offset int) (models.GoodsResponse, error) {
	goodsResponse, err := s.repoPostgres.Goods(ctx, limit, offset)
	if err != nil {
		return models.GoodsResponse{}, fmt.Errorf("error to get goods: %v", err)
	}
	return goodsResponse, nil
}

func (s *Service) CreateGood(ctx context.Context, projectID int, name string) (models.Good, error) {
	good, err := s.repoPostgres.CreateGood(ctx, projectID, name)
	if err != nil {
		return models.Good{}, fmt.Errorf("error to create good: %v", err)
	}
	return good, nil
}

func (s *Service) UpdateGood(ctx context.Context, good models.Good, goodID, projectID int) (models.Good, error) {
	good, err := s.repoPostgres.UpdateGood(ctx, good, goodID, projectID)
	if err != nil {
		return models.Good{}, fmt.Errorf("error to update good: %v", err)
	}
	return good, nil
}