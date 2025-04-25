package service

import (
	"backend/internal/application/interface"
	model2 "backend/internal/application/model"
	"backend/internal/domain/model"
	"backend/internal/domain/repository"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"math/rand"
	"sync"
	"time"
)

var _ _interface.ProductService = (*ProductServiceImpl)(nil)

type ProductServiceImpl struct {
	productRepository repository.ProductRepository
	cacheService      _interface.CacheService
}

func NewProductServiceImpl(
	productRepository repository.ProductRepository,
	cacheService _interface.CacheService,
) *ProductServiceImpl {
	return &ProductServiceImpl{
		productRepository: productRepository,
		cacheService:      cacheService,
	}
}

const (
	DashboardCacheKey         = "dashboard_information"
	BatchSize                 = 1000
	TotalRecords              = 10000000
	TestDataGenerationLockKey = "test_data_generation_lock"
	NumWorkers                = 20
)

func (p *ProductServiceImpl) GetProducts(ctx context.Context, filter *repository.Filter) ([]model2.GetProductsResponse, int64, error) {
	filter.Offset = filter.GetOffset()

	products, totalRows, err := p.productRepository.Get(ctx, filter)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "Failed to get products")
	}

	var productResponses []model2.GetProductsResponse
	for _, product := range products {
		productResponses = append(productResponses, model2.ToGetProductResponse(product))
	}
	return productResponses, totalRows, nil
}

func (p *ProductServiceImpl) Add(ctx context.Context, request model2.AddProductRequest) error {
	status := model.InStock
	stockQuantity := request.StockQuantity
	if stockQuantity <= 0 {
		status = model.OutOfStock
	}
	products := []model.Product{
		{
			Name:          request.Name,
			Description:   request.Description,
			Price:         request.Price,
			StockQuantity: stockQuantity,
			Status:        string(status),
			CategoryId:    request.CategoryId,
		},
	}

	if err := p.productRepository.SaveAll(ctx, products); err != nil {
		return errors.Wrapf(err, "Failed to save product")
	}

	p.clearDashboardCache()
	return nil
}

func (p *ProductServiceImpl) GetDashboard(ctx context.Context) ([]model2.GetDashboardResponse, error) {
	var response []model2.GetDashboardResponse
	err := p.cacheService.Get(ctx, DashboardCacheKey, &response)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get dashboard cache")
	} else {
		log.Info().Msg("Cache hit for dashboard information")
		return response, nil
	}

	dashboard, err := p.productRepository.GetDashboard(ctx)
	if err != nil {
		return nil, err
	}

	for _, item := range dashboard {
		response = append(response, model2.GetDashboardResponse{
			CategoryName: item.CategoryName,
			ProductCount: item.ProductCount,
		})
	}

	go func() {
		err = p.cacheService.Set(context.TODO(), DashboardCacheKey, response, 1*time.Hour)
		if err != nil {
			log.Error().Err(err).Msg("Failed to set dashboard cache")
		}
	}()
	return response, nil
}

func (p *ProductServiceImpl) GenerateTestData(ctx context.Context) error {
	startTime := time.Now()
	log.Info().Msgf("Starting test data generation at %s", startTime.Format(time.RFC3339))
	lock, err := p.cacheService.ObtainLock(ctx, TestDataGenerationLockKey, 5*time.Minute)
	if err != nil {
		return errors.Wrapf(err, "Another test data generation process is running")
	}
	defer func() {
		p.clearDashboardCache()
		log.Info().Msgf("Took %s to generate test data", time.Since(startTime).String())
		lock.Release(ctx)
	}()

	productChan := make(chan model.Product, BatchSize)
	var wg sync.WaitGroup
	errChan := make(chan error, NumWorkers)

	// Worker function
	worker := func(workerId int64) {
		defer wg.Done()
		var products []model.Product

		for product := range productChan {
			products = append(products, product)

			// When batch size is reached, save and clear the slice
			if len(products) == BatchSize {
				// to csv
				err := p.productRepository.SaveAll(ctx, products)
				if err != nil {
					errChan <- err
					return
				}
				log.Info().Msgf("Worker %d saved batch of %d products", workerId, len(products))
				products = nil // Clear the slice
			}
		}

		// Save any remaining products
		if len(products) > 0 {
			log.Info().Msgf("Worker %d saved %d products", workerId, len(products))
			err := p.productRepository.SaveAll(ctx, products)
			if err != nil {
				errChan <- err
			}
		}
	}

	// Start workers
	for i := 1; i <= NumWorkers; i++ {
		wg.Add(1)
		go worker(int64(i))
	}

	// Generate products and send to workers
	go func() {
		for i := 1; i <= TotalRecords; i++ {
			productChan <- model.Product{
				Name:          fmt.Sprintf("Product %d", i),
				Description:   fmt.Sprintf("Description %d", i),
				Price:         float64(rand.Intn(10000)) / 100,
				StockQuantity: int64(rand.Intn(10000)),
				Status:        "in_stock",
				CategoryId:    1,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			}
		}
		close(productChan)
	}()

	// Wait for workers to finish
	wg.Wait()
	close(errChan)

	// Check for errors
	for err := range errChan {
		log.Error().Err(err).Msg("Batch processing error")
	}
	log.Info().Msg("Test data generation completed successfully")
	return nil
}

// Clear the cache after add/edit a new product
func (p *ProductServiceImpl) clearDashboardCache() {
	err := p.cacheService.Del(context.TODO(), DashboardCacheKey)
	if err != nil {
		log.Error().Err(err).Msg("Failed to clear dashboard cache")
		return
	}
	log.Info().Msg("Cleared dashboard cache successfully")
}
