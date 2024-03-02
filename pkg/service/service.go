package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Hymiside/hezzl-api/pkg/custerrors"
	"github.com/Hymiside/hezzl-api/pkg/models"
	log "github.com/sirupsen/logrus"
)

type repositoryPostgres interface{
	Goods(ctx context.Context) ([]models.Good, error)
	GoodsWithLimitAndOffset(ctx context.Context, limit, offset int) (models.GoodsResponse, error)
	CreateGood(ctx context.Context, projectID int, name string) (models.Good, error)
	UpdateGood(ctx context.Context, good models.Good, goodID, projectID int) (models.Good, error)
	DeleteGood(ctx context.Context, goodID, projectID int) (models.Good, error)
	ReprioritizeGood(ctx context.Context, goodID, projectID, priority int) ([]models.ReprioritizeGoodResponse, error)
}

type repositoryRedis interface{
	Set(ctx context.Context, goods []byte) error
	Get(ctx context.Context, limit, offset int) (models.GoodsResponse, error)
	Delete(ctx context.Context) error
}

type queueNats interface{
	Subscribe() error
	Publish(b []byte) error
}

type Service struct {
	repoPostgres   repositoryPostgres
	repoRedis      repositoryRedis
	queueNats      queueNats
}

func NewService(
	repoPostgres repositoryPostgres,
	repoRedis repositoryRedis,
	queueNats queueNats,
) *Service {

	if err := queueNats.Subscribe(); err != nil {
		log.Fatalf("error to subscribe queue nats: %v", err)
	}

	return &Service{
		repoPostgres:   repoPostgres,
		repoRedis:      repoRedis,
		queueNats:      queueNats,
	}
}

func (s *Service) cacheGoods(ctx context.Context) {
	goods, err := s.repoPostgres.Goods(ctx)
	if err != nil {
		log.Errorf("error to get goods: %v", err)
	}

	bytes, err := json.Marshal(goods)
	if err != nil {
		log.Errorf("error to marshal goods: %v", err)
	}
	
	if err := s.repoRedis.Set(ctx, bytes); err != nil {
		log.Errorf("error to set goods: %v", err)
	}
}

func (s *Service) Goods(ctx context.Context, limit, offset int) (models.GoodsResponse, error) {
	responseGoods, err := s.repoRedis.Get(ctx, limit, offset)
	if err != nil {
		if !errors.Is(err, custerrors.ErrNotFound) {
			return models.GoodsResponse{}, fmt.Errorf("error to get goods: %v", err)
		}

		responseGoods, err = s.repoPostgres.GoodsWithLimitAndOffset(ctx, limit, offset)
		if err != nil {
			return models.GoodsResponse{}, fmt.Errorf("error to get goods: %v", err)
		}
		
		go s.cacheGoods(context.Background())
	}
	
	return responseGoods, nil
}

func (s *Service) CreateGood(ctx context.Context, projectID int, name string) (models.Good, error) {
	createdGood, err := s.repoPostgres.CreateGood(ctx, projectID, name)
	if err != nil {
		return models.Good{}, fmt.Errorf("error to create good: %v", err)
	}

	bytes, err := json.Marshal(createdGood)
	if err != nil {
		log.Errorf("error to marshal good: %v", err)
	}

	if err = s.queueNats.Publish(bytes); err != nil {
		log.Errorf("error to publish good: %v", err)
	}

	if err = s.repoRedis.Delete(ctx); err != nil {
		log.Errorf("error to delete redis: %v", err)
	}
	return createdGood, nil
}

func (s *Service) UpdateGood(ctx context.Context, good models.Good, goodID, projectID int) (models.Good, error) {
	updatedGood, err := s.repoPostgres.UpdateGood(ctx, good, goodID, projectID)
	if err != nil {
		if errors.Is(err, custerrors.ErrNotFound) {
			return models.Good{}, err
		}
		return models.Good{}, fmt.Errorf("error to update good: %v", err)
	}

	bytes, err := json.Marshal(updatedGood)
	if err != nil {
		log.Errorf("error to marshal good: %v", err)
	}

	if err = s.queueNats.Publish(bytes); err != nil {
		log.Errorf("error to publish good: %v", err)
	}

	if err = s.repoRedis.Delete(ctx); err != nil {
		log.Errorf("error to delete redis: %v", err)
	}

	return updatedGood, nil
}

func (s *Service) DeleteGood(ctx context.Context, goodID, projectID int) (models.Good, error) {
	deletedGood, err := s.repoPostgres.DeleteGood(ctx, goodID, projectID)
	if err != nil {
		if errors.Is(err, custerrors.ErrNotFound) {
			return models.Good{}, err
		}
		return models.Good{}, fmt.Errorf("error to delete good: %v", err)
	}

	bytes, err := json.Marshal(deletedGood)
	if err != nil {
		log.Errorf("error to marshal good: %v", err)
	}

	if err = s.queueNats.Publish(bytes); err != nil {
		log.Errorf("error to publish good: %v", err)
	}

	if err = s.repoRedis.Delete(ctx); err != nil {
		log.Errorf("error to delete redis: %v", err)
	}
	return deletedGood, nil
}

func (s *Service) ReprioritizeGood(ctx context.Context, goodID, projectID, priority int) ([]models.ReprioritizeGoodResponse, error) {
	reprioritizedGoods, err := s.repoPostgres.ReprioritizeGood(ctx, goodID, projectID, priority)
	if err != nil {
		if errors.Is(err, custerrors.ErrNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("error to reprioritize good: %v", err)
	}

	bytes, err := json.Marshal(reprioritizedGoods)
	if err != nil {
		log.Errorf("error to marshal goods: %v", err)
	}

	if err = s.queueNats.Publish(bytes); err != nil {
		log.Errorf("error to publish good: %v", err)
	}

	if err = s.repoRedis.Delete(ctx); err != nil {
		log.Errorf("error to delete redis: %v", err)
	}
	
	return reprioritizedGoods, nil
}