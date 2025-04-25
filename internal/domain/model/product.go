package model

import "time"

type Product struct {
	Id            int64
	Price         float64
	Name          string
	Description   string
	StockQuantity int64
	Status        string
	CategoryId    int64
	CategoryName  string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Dashboard struct {
	CategoryName string
	ProductCount int64
}

type ProductStatus string

const (
	InStock    ProductStatus = "in_stock"
	OutOfStock ProductStatus = "out_of_stock"
)
