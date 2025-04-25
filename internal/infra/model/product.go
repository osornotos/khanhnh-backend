package model

import (
	"github.com/uptrace/bun"
	"time"
)

type ProductEntity struct {
	bun.BaseModel `bun:"table:products"`

	Id            int64   `bun:"id,pk,autoincrement"`
	Price         float64 `bun:"price"`
	Name          string  `bun:"name"`
	Description   string  `bun:"description"`
	StockQuantity int64   `bun:"stock_quantity"`
	CategoryId    int64   `bun:"category_id"`
	Status        string  `bun:"status"`

	CreatedAt time.Time `bun:"created_at"`
	UpdatedAt time.Time `bun:"updated_at"`
}

type DashboardProductCount struct {
	CategoryName string `bun:"category_name"`
	ProductCount int    `bun:"product_count"`
}
