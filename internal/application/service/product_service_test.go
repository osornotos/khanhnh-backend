package service

import (
	"context"
	"khanhnh-backend/internal/domain/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	model2 "khanhnh-backend/internal/application/model"
	"khanhnh-backend/internal/domain/model"
	"khanhnh-backend/tests/mocks"
)

func TestProductServiceImpl_Get(t *testing.T) {
	mockRepo := new(mocks.ProductRepositoryMock)
	mockCacheService := new(mocks.CacheServiceMock)
	service := NewProductServiceImpl(mockRepo, mockCacheService)

	ctx := context.Background()
	filter := &repository.Filter{
		PageRequest: repository.PageRequest{}}
	mockProducts := []model.Product{
		{Name: "Product1", Description: "Desc1", Price: 100, StockQuantity: 10, Status: "InStock"},
	}
	mockRepo.On("Get", ctx, filter).Return(mockProducts, int64(1), nil)

	products, totalRows, err := service.GetProducts(ctx, filter)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), totalRows)
	assert.Len(t, products, 1)
	assert.Equal(t, "Product1", products[0].Name)
	mockRepo.AssertExpectations(t)
}

func TestProductServiceImpl_Add(t *testing.T) {
	mockRepo := new(mocks.ProductRepositoryMock)
	mockCacheService := new(mocks.CacheServiceMock)
	service := NewProductServiceImpl(mockRepo, mockCacheService)

	ctx := context.Background()
	request := model2.AddProductRequest{
		Name:          "Product1",
		Description:   "Desc1",
		Price:         100,
		StockQuantity: 10,
	}
	expectedProducts := []model.Product{
		{
			Name:          "Product1",
			Description:   "Desc1",
			Price:         100,
			StockQuantity: 10,
			Status:        "in_stock",
		},
	}
	mockRepo.On("SaveAll", ctx, expectedProducts).Return(nil)
	mockCacheService.On("Del", context.TODO(), "dashboard_information").Return(nil)

	err := service.Add(ctx, request)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
