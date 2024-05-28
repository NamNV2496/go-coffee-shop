package cache

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/namnv2496/go-coffee-shop-demo/internal/configs"
	"github.com/namnv2496/go-coffee-shop-demo/internal/mq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrCacheMiss = errors.New("cache miss")
)

type Client interface {
	Set(ctx context.Context, key string, data any, ttl ...time.Duration) error
	Get(ctx context.Context, key string) (any, error)
	AddToSet(ctx context.Context, key string, data ...any) error
	IsDataInSet(ctx context.Context, key string, data any) (bool, error)
}

type redisClient struct {
	redisClient *redis.Client
}

func NewRedisClient(
	cacheConfig configs.Redis,
) Client {
	return &redisClient{
		redisClient: redis.NewClient(&redis.Options{
			Addr:     cacheConfig.Address,
			Username: cacheConfig.Username,
			Password: cacheConfig.Password,
			DB:       cacheConfig.Database,
		}),
	}
}

func (c redisClient) Set(ctx context.Context, key string, data any, ttl ...time.Duration) error {

	var ttlValue time.Duration
	if len(ttl) > 0 {
		ttlValue = ttl[0]
	} else {
		ttlValue = mq.RedisTTL
	}
	if err := c.redisClient.Set(ctx, key, data, ttlValue).Err(); err != nil {
		return status.Error(codes.Internal, "failed to set data into cache")
	}

	return nil
}

func (c redisClient) Get(ctx context.Context, key string) (any, error) {

	data, err := c.redisClient.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrCacheMiss
		}
		return nil, status.Error(codes.Internal, "failed to get data from cache")
	}

	return data, nil
}

func (c redisClient) AddToSet(ctx context.Context, key string, data ...any) error {

	if err := c.redisClient.SAdd(ctx, key, data...).Err(); err != nil {
		return status.Error(codes.Internal, "failed to set data into set inside cache")
	}

	return nil
}

func (c redisClient) IsDataInSet(ctx context.Context, key string, data any) (bool, error) {

	result, err := c.redisClient.SIsMember(ctx, key, data).Result()
	if err != nil {
		return false, status.Error(codes.Internal, "failed to check if data is member of set inside cache")
	}

	return result, nil
}

type inMemoryClient struct {
	cache      map[string]any
	cacheMutex *sync.Mutex
}
