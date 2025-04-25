package model

type AddProductRequest struct {
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Price         float64 `json:"price"`
	StockQuantity int64   `json:"stockQuantity"`
	CategoryId    int64   `json:"categoryId"`
}
