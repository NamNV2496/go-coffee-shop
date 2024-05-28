package service

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/namnv2496/go-coffee-shop-demo/internal/cache"
	"github.com/namnv2496/go-coffee-shop-demo/internal/mq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type KitchenService interface {
	UpdateStatusOrderToRedis(ctx context.Context, customerId int32, itemId int32, finished int32) error
	GetOrderInRedis(ctx context.Context) ([]mq.RedisOrderDTO, error)
	GetCustomerOrderInRedis(ctx context.Context, customerId int32) (mq.RedisOrderDTO, error)
}

type kitchenService struct {
	Cache cache.Client
}

func NewService(
	cache cache.Client,
) KitchenService {

	return &kitchenService{
		Cache: cache,
	}
}

func (k kitchenService) UpdateStatusOrderToRedis(
	ctx context.Context,
	customerId int32,
	itemId int32,
	finished int32,
) error {

	data, exist := k.Cache.Get(ctx, mq.REDIS_KEY_ORDER)
	if exist != nil {
		return status.Error(codes.Internal, "failed to update status of not exist customerId")
	}
	var jsonData []mq.RedisOrderDTO
	publicKey, _ := data.(string)
	if err := json.Unmarshal([]byte(publicKey), &jsonData); err != nil {
		return status.Error(codes.Internal, "failed to unmarshal data")
	}
	for i, customerOrder := range jsonData {
		if customerOrder.CustomerId == customerId {
			var redisOrderNew = make([]mq.RedisOrder, 0)
			for _, orderItem := range customerOrder.RedisOrders {
				if orderItem.ItemId == itemId {
					if orderItem.Status == mq.Canceled {
						return status.Error(codes.Internal, "Cannot add finish for cancel order")
					}
					newFinished := orderItem.Finished + finished
					orderItem.Finished = newFinished
					if newFinished == orderItem.Quantity {
						orderItem.Status = mq.Done
					} else if newFinished > orderItem.Quantity {
						return status.Error(codes.Internal, "Cannot add finished value > order quantity")
					}
				}
				redisOrderNew = append(redisOrderNew, orderItem)
			}
			jsonData[i].RedisOrders = redisOrderNew
		}
	}
	json, ok := json.Marshal(jsonData)
	if ok != nil {
		return errors.New("failed to marshall data into cache")
	}
	k.Cache.Set(ctx, mq.REDIS_KEY_ORDER, json)
	return nil
}

func (k kitchenService) GetOrderInRedis(ctx context.Context) ([]mq.RedisOrderDTO, error) {
	data, exist := k.Cache.Get(ctx, mq.REDIS_KEY_ORDER)
	if exist != nil {
		return nil, status.Error(codes.Internal, "failed to update status of not exist customerId")
	}
	var jsonData []mq.RedisOrderDTO
	publicKey, _ := data.(string)
	if err := json.Unmarshal([]byte(publicKey), &jsonData); err != nil {
		return nil, status.Error(codes.Internal, "failed to unmarshal data")
	}
	return jsonData, nil
}

func (k kitchenService) GetCustomerOrderInRedis(ctx context.Context, customerId int32) (mq.RedisOrderDTO, error) {
	data, exist := k.Cache.Get(ctx, mq.REDIS_KEY_ORDER)
	if exist != nil {
		return mq.RedisOrderDTO{}, status.Error(codes.Internal, "failed to update status of not exist customerId")
	}
	var jsonData []mq.RedisOrderDTO
	publicKey, _ := data.(string)
	if err := json.Unmarshal([]byte(publicKey), &jsonData); err != nil {
		return mq.RedisOrderDTO{}, status.Error(codes.Internal, "failed to unmarshal data")
	}

	for _, redisOrderDTO := range jsonData {
		if redisOrderDTO.CustomerId == customerId {
			return redisOrderDTO, nil
		}
	}
	return mq.RedisOrderDTO{}, nil
}
