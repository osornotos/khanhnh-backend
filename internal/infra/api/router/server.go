package router

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/rs/zerolog/log"
	"khanhnh-backend/internal"
	"khanhnh-backend/wire"
	"os"
	"os/signal"
	"syscall"
)

func StartServer(app *wire.Application) {
	cfg := app.Config
	_, cancelFunc := context.WithCancel(context.Background())

	r := NewRouter(cfg)

	r.WireRoot(cfg)
	r.WireProductRoutes(app)
	r.WireUserRoutes(app)
	go r.Start(cfg.Port)

	interruptChan := make(chan os.Signal)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interruptChan

	r.Stop()
	cancelFunc()
	log.Info().Msgf("shutdown http server successfully!")
}

func setupMiddleware(app *fiber.App, cfg internal.Config) {
	app.Use(requestid.New())
	app.Use(helmet.New())
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))
	setUpLoggerMiddleware(app)
}

func setUpLoggerMiddleware(app *fiber.App) {
	logConfig := logger.Config{
		Format:     `{"time":"${time}","status":"${status}","method":"${method}","url":"${url}","latency":"${latency}","remote_ip":"${ip}","referer":"${referer}","host":"${host}","user_agent":"${ua}","bytes_in":"${bytesReceived}","bytes_out":"${bytesSent}","err":"${error}","request_id":"${locals:requestid}}"` + "\n",
		TimeFormat: "2006-01-02T15:04:05.0000-0700",
		TimeZone:   "UTC",
		Output:     os.Stdout,
		CustomTags: map[string]logger.LogFunc{
			"cache": func(output logger.Buffer, c *fiber.Ctx, data *logger.Data, extraParam string) (int, error) {
				cache := c.Get("cache")
				return output.WriteString(cache)
			},
		},
	}
	app.Use(logger.New(logConfig))
}
