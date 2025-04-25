package controller

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"khanhnh-backend/internal/application/interface"
	"khanhnh-backend/internal/application/model"
	"khanhnh-backend/internal/domain/repository"
	"net/http"
)

type ProductController struct {
	productService _interface.ProductService
}

func NewProductController(productService _interface.ProductService) *ProductController {
	return &ProductController{
		productService: productService,
	}
}

func (w *ProductController) Get(ctx *fiber.Ctx) error {
	params, err := parseRequestQuery[repository.Filter](ctx)
	products, totalRows, err := w.productService.GetProducts(ctx.Context(), &params)
	if err != nil {
		log.Error().Stack().Err(err).Msg("get products")
		return err
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"data":   products,
		"status": "ok",
		"total":  totalRows,
		"page":   params.Page,
	})
}

func (w *ProductController) Add(ctx *fiber.Ctx) error {
	params, err := parseRequestParams[model.AddProductRequest](ctx)
	if err != nil {
		return err
	}
	err = w.productService.Add(ctx.Context(), params)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"status": "ok",
	})
}

func (w *ProductController) GetDashboard(ctx *fiber.Ctx) error {
	dashboard, err := w.productService.GetDashboard(ctx.Context())
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"data":   dashboard,
		"status": "ok",
	})
}

func (w *ProductController) GenerateTestData(ctx *fiber.Ctx) error {
	go func() {
		err := w.productService.GenerateTestData(context.Background())
		if err != nil {
			log.Error().Stack().Err(err).Msg("Generate test data")
		}
	}()

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"status": "ok",
	})
}
