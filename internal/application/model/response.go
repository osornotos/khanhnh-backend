package model

import (
	"khanhnh-backend/internal/domain/model"
	"time"
)

type GetProductsResponse struct {
	Id            int64   `json:"id"`
	Price         float64 `json:"price"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	StockQuantity int64   `json:"stockQuantity"`
	Status        string  `json:"status"`
	CategoryName  string  `json:"categoryName"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func ToGetProductResponse(product model.Product) GetProductsResponse {
	return GetProductsResponse{
		Id:            product.Id,
		Price:         product.Price,
		Name:          product.Name,
		Description:   product.Description,
		StockQuantity: product.StockQuantity,
		Status:        product.Status,
		CategoryName:  product.CategoryName,
		CreatedAt:     product.CreatedAt,
		UpdatedAt:     product.UpdatedAt,
	}
}

type GetDashboardResponse struct {
	CategoryName string `json:"categoryName"`
	ProductCount int64  `json:"productCount"`
}
