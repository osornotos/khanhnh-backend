package repository

import (
	"context"
	"khanhnh-backend/internal/domain/model"
)

type Filter struct {
	Id          int64
	IdIn        []int64
	Description string
	Name        string
	Statuses    []string
	PageRequest
}

type PageRequest struct {
	Page     int64
	PageSize int64
	Offset   int64
}

const DefaultPageSize = 20

func (p *PageRequest) GetOffset() int64 {
	if p.Page == 0 {
		p.Page = 1
	}
	if p.PageSize == 0 {
		p.PageSize = DefaultPageSize
	}
	return (p.Page - 1) * p.PageSize
}

type ProductRepository interface {
	Get(ctx context.Context, filter *Filter) ([]model.Product, int64, error)
	SaveAll(ctx context.Context, products []model.Product) error
	GetDashboard(ctx context.Context) ([]model.Dashboard, error)
}
