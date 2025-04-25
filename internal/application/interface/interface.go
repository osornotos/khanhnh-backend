package _interface

import (
	"backend/internal/application/model"
	"backend/internal/domain/repository"
	"context"
	"time"
)

type ProductService interface {
	GetProducts(ctx context.Context, filter *repository.Filter) ([]model.GetProductsResponse, int64, error)
	Add(ctx context.Context, request model.AddProductRequest) error
	GetDashboard(ctx context.Context) ([]model.GetDashboardResponse, error)
	GenerateTestData(ctx context.Context) error
}

type UserService interface {
	Login(ctx context.Context) (string, error)
}

type CacheService interface {
	Get(ctx context.Context, key string, result interface{}) (err error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Del(ctx context.Context, key string) error
	ObtainLock(ctx context.Context, key string, duration time.Duration) (Locker, error)
}

type Locker interface {
	Release(ctx context.Context) error
}
