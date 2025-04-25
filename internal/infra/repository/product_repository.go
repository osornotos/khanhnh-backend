package repository

import (
	"backend/internal/domain/model"
	"backend/internal/domain/repository"
	model2 "backend/internal/infra/model"
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
	"time"
)

var _ repository.ProductRepository = (*ProductRepositoryImpl)(nil)

type ProductRepositoryImpl struct {
	db *bun.DB
}

const ProductsTable = "products"
const CategoriesTable = "categories"

func NewProductRepositoryImpl(db *bun.DB) *ProductRepositoryImpl {
	return &ProductRepositoryImpl{
		db: db,
	}
}

func bindFormFilter[T QueryCondition[T]](query T, filter *repository.Filter) T {
	if filter.Name != "" {
		query.Where("products.name ILIKE ?", sqlLike(filter.Name))
	}
	if filter.Description != "" {
		query.Where("products.description ILIKE ?", sqlLike(filter.Description))
	}
	if len(filter.Statuses) > 0 {
		query.Where("products.status IN (?)", bun.In(filter.Statuses))
	}
	return query
}

//func (p *ProductRepositoryImpl) Get(ctx context.Context, filter *repository.Filter) ([]model.Product, int64, error) {
//	var productEntities []model2.ProductEntity
//	queryBuilder := p.db.NewSelect().Model(&productEntities)
//	queryBuilder = bindFormFilter(queryBuilder, filter)
//
//	sortBy := "created_at DESC"
//
//	totalRows, err := queryBuilder.
//		OrderExpr(sortBy).
//		Offset(int(filter.Offset)).
//		Limit(int(filter.PageSize)).
//		ScanAndCount(ctx)
//	if err != nil && !errors.Is(err, sql.ErrNoRows) {
//		return nil, 0, errors.WithStack(err)
//	}
//	if errors.Is(err, sql.ErrNoRows) {
//		return []model.Product{}, 0, nil
//	}
//	var products []model.Product
//	for _, productEntity := range productEntities {
//		products = append(products, toProduct(productEntity))
//	}
//	return products, int64(totalRows), nil
//}

func (p *ProductRepositoryImpl) Get(ctx context.Context, filter *repository.Filter) ([]model.Product, int64, error) {
	var productEntities []struct {
		model2.ProductEntity
		CategoryId          int64  `bun:"category_id"`
		CategoryName        string `bun:"category_name"`
		CategoryDescription string `bun:"category_description"`
	}
	queryBuilder := p.db.NewSelect().
		Table("products").
		Model(&productEntities).
		ColumnExpr("products.*").
		ColumnExpr("categories.id AS category_id, categories.name AS category_name, categories.description AS category_description").
		Join("LEFT JOIN categories ON products.category_id = categories.id")

	queryBuilder = bindFormFilter(queryBuilder, filter)

	sortBy := "products.created_at DESC"

	totalRows, err := queryBuilder.
		OrderExpr(sortBy).
		Offset(int(filter.Offset)).
		Limit(int(filter.PageSize)).
		ScanAndCount(ctx)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, 0, errors.WithStack(err)
	}
	if errors.Is(err, sql.ErrNoRows) {
		return []model.Product{}, 0, nil
	}
	var products []model.Product
	for _, entity := range productEntities {
		product := toProduct(entity.ProductEntity)
		product.CategoryId = entity.CategoryId
		product.CategoryName = entity.CategoryName
		product.Description = entity.CategoryDescription
		products = append(products, product)
	}
	return products, int64(totalRows), nil
}

func (p *ProductRepositoryImpl) SaveAll(ctx context.Context, products []model.Product) error {
	var productEntities []model2.ProductEntity
	for _, product := range products {
		if product.CreatedAt.IsZero() {
			product.CreatedAt = time.Now().UTC()
		}
		if product.UpdatedAt.IsZero() {
			product.CreatedAt = time.Now().UTC()
		}
		productEntities = append(productEntities, toProductEntity(product))
	}
	_, err := p.db.NewInsert().
		Model(&productEntities).
		Exec(ctx)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func toProductEntity(product model.Product) model2.ProductEntity {
	return model2.ProductEntity{
		Id:            product.Id,
		Name:          product.Name,
		Description:   product.Description,
		Price:         product.Price,
		StockQuantity: product.StockQuantity,
		Status:        product.Status,
		CategoryId:    product.CategoryId,
		CreatedAt:     product.CreatedAt,
		UpdatedAt:     product.UpdatedAt,
	}
}

func toProduct(productEntity model2.ProductEntity) model.Product {
	return model.Product{
		Id:            productEntity.Id,
		Name:          productEntity.Name,
		Description:   productEntity.Description,
		Price:         productEntity.Price,
		StockQuantity: productEntity.StockQuantity,
		Status:        productEntity.Status,
		CreatedAt:     productEntity.CreatedAt,
		UpdatedAt:     productEntity.UpdatedAt,
	}
}

func (p *ProductRepositoryImpl) GetDashboard(ctx context.Context) ([]model.Dashboard, error) {
	var results []model2.DashboardProductCount

	err := p.db.NewSelect().
		Table("categories").
		ColumnExpr("categories.name AS category_name").
		ColumnExpr("COUNT(products.id) AS product_count").
		Join("LEFT JOIN products ON categories.id = products.category_id").
		Group("categories.id").
		Scan(ctx, &results)

	if err != nil {
		return nil, err
	}
	var dashboard []model.Dashboard
	for _, result := range results {
		dashboard = append(dashboard, model.Dashboard{
			CategoryName: result.CategoryName,
			ProductCount: int64(result.ProductCount),
		})
	}
	return dashboard, nil
}
