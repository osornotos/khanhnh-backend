package service

import (
	"backend/internal/application/interface"
	"encoding/json"
	"fmt"
	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
	"time"

	"context"
)

type CacheServiceImpl struct {
	client *redis.Client
}

var _ _interface.CacheService = (*CacheServiceImpl)(nil)

func NewCacheServiceImpl(client *redis.Client) *CacheServiceImpl {
	return &CacheServiceImpl{
		client: client,
	}
}

const BackendPrefix = "backend"

func (c *CacheServiceImpl) GetKey(key string) string {
	return fmt.Sprintf("%v:%v", BackendPrefix, key)
}

func (c *CacheServiceImpl) Get(ctx context.Context, key string, result interface{}) error {
	key = c.GetKey(key)
	rawData, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(rawData), result)
	if err != nil {
		return err
	}
	return nil
}

func (c *CacheServiceImpl) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	key = c.GetKey(key)
	rawData, _ := json.Marshal(value)
	return c.client.Set(ctx, key, rawData, expiration).Err()
}

func (c *CacheServiceImpl) Del(ctx context.Context, key string) error {
	key = c.GetKey(key)
	return c.client.Del(ctx, key).Err()
}

func (c *CacheServiceImpl) ObtainLock(ctx context.Context, key string, duration time.Duration) (_interface.Locker, error) {
	locker := redislock.New(c.client)
	key = c.GetKey(key)
	retryOptions := redislock.Options{RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(300*time.Millisecond), 3)}
	lock, err := locker.Obtain(ctx, key, duration, &retryOptions)
	if err != nil {
		return nil, err
	}
	return lock, nil
}
