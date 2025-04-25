package service

import (
	"context"
	"github.com/alicebob/miniredis/v2"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func setupTestCacheService() (*CacheServiceImpl, func()) {
	mr, _ := miniredis.Run()
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	cacheService := NewCacheServiceImpl(client)
	return cacheService, func() { mr.Close() }
}

func TestCacheServiceImpl_GetKey(t *testing.T) {
	cacheService, teardown := setupTestCacheService()
	defer teardown()

	key := "test-key"
	expected := "backend:test-key"
	assert.Equal(t, expected, cacheService.GetKey(key))
}

func TestCacheServiceImpl_SetAndGet(t *testing.T) {
	cacheService, teardown := setupTestCacheService()
	defer teardown()

	ctx := context.Background()
	key := "test-key"
	value := map[string]string{"field": "value"}

	err := cacheService.Set(ctx, key, value, time.Minute)
	assert.NoError(t, err)

	var result map[string]string
	err = cacheService.Get(ctx, key, &result)
	assert.NoError(t, err)
	assert.Equal(t, value, result)
}

func TestCacheServiceImpl_Del(t *testing.T) {
	cacheService, teardown := setupTestCacheService()
	defer teardown()

	ctx := context.Background()
	key := "test-key"
	value := "test-value"

	err := cacheService.Set(ctx, key, value, time.Minute)
	assert.NoError(t, err)

	err = cacheService.Del(ctx, key)
	assert.NoError(t, err)

	var result string
	err = cacheService.Get(ctx, key, &result)
	assert.Error(t, err)
}

func TestCacheServiceImpl_ObtainLock(t *testing.T) {
	cacheService, teardown := setupTestCacheService()
	defer teardown()

	ctx := context.Background()
	key := "test-lock"
	duration := time.Second * 5

	lock, err := cacheService.ObtainLock(ctx, key, duration)
	assert.NoError(t, err)
	assert.NotNil(t, lock)

	// Try to obtain the same lock again
	_, err = cacheService.ObtainLock(ctx, key, duration)
	assert.Error(t, err)

	// Release the lock
	err = lock.Release(ctx)
	assert.NoError(t, err)

	// Try to obtain the lock again
	lock, err = cacheService.ObtainLock(ctx, key, duration)
	assert.NoError(t, err)
}
