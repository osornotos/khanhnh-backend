package router

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"khanhnh-backend/internal"
	"khanhnh-backend/internal/infra/api/controller"
	"khanhnh-backend/internal/infra/middleware"
	"khanhnh-backend/wire"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type Router struct {
	app *fiber.App
}

func NewRouter(cfg internal.Config) *Router {
	httpServer := fiber.New(fiber.Config{
		CaseSensitive:           true,
		ProxyHeader:             fiber.HeaderXForwardedFor,
		EnableTrustedProxyCheck: false,
		ReadTimeout:             10 * time.Second,
		BodyLimit:               11 * 1024 * 1024,

		JSONEncoder:  json.Marshal,
		ErrorHandler: middleware.ErrorHandler,
	})

	setupMiddleware(httpServer, cfg)
	r := &Router{app: httpServer}
	return r
}

func (r *Router) Start(port int) {
	if err := r.app.Listen(fmt.Sprintf(":%v", port)); err != nil && err != http.ErrServerClosed {
		log.Error().Stack().Err(err).Msg("HTTP server failed to start")
	}
}

func (r *Router) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if err := r.app.ShutdownWithContext(ctx); err != nil {
		log.Fatal().Err(err).Msgf("failed to shutdown http server gracefully")
	}
}

func (r *Router) WireRoot(cfg internal.Config) {
	r.app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Status(http.StatusOK).JSON(fiber.Map{
			"app":     "khanhnh-backend",
			"version": cfg.Version,
		})
	})
}

func (r *Router) WireProductRoutes(app *wire.Application) {
	tokenRequired := middleware.Authorization(app.JwtProvider)
	//service := r.app.Group("service", tokenRequired)
	product := r.app.Group("products", tokenRequired)
	productController := controller.NewProductController(app.ProductService)
	product.Get("", productController.Get)
	product.Get("dashboard", productController.GetDashboard)
	product.Post("generate-test-data", productController.GenerateTestData)
}

func (r *Router) WireUserRoutes(app *wire.Application) {
	user := r.app.Group("users")
	userController := controller.NewUserController(app.UserService)
	user.Post("login", userController.Login)
}

func (r *Router) GetApp() *fiber.App {
	return r.app
}
