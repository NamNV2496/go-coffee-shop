package service

import (
	"context"
	"encoding/json"

	"github.com/namnv2496/go-coffee-shop-demo/internal/batch/domain"
	"github.com/namnv2496/go-coffee-shop-demo/internal/batch/repo"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/cache"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/mq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BatchService interface {
	ClearAllOrderEOD(ctx context.Context) error
	GetOrderOfToday(ctx context.Context) ([]domain.OrderItem, error)
}

type batchService struct {
	orderRepo repo.OrderRepo
	cache     cache.Client
}

func NewBatchService(
	orderRepo repo.OrderRepo,
	cache cache.Client,
) BatchService {
	return &batchService{
		orderRepo: orderRepo,
		cache:     cache,
	}
}

func (s batchService) ClearAllOrderEOD(
	ctx context.Context,
) error {

	var jsonData = make([]mq.RedisOrderDTO, 0)
	json, ok := json.Marshal(jsonData)
	if ok != nil {
		return status.Error(codes.Internal, "failed to marshall data into cache")
	}
	s.cache.Set(ctx, mq.REDIS_KEY_ORDER, json)
	return nil
}

func (s batchService) GetOrderOfToday(ctx context.Context) ([]domain.OrderItem, error) {

	orderItems, err := s.orderRepo.GetOrderItems(ctx)
	if err != nil {
		return []domain.OrderItem{}, nil
	}
	return orderItems, nil
}
