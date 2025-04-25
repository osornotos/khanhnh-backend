package wire

import (
	"backend/internal"
	"backend/internal/application/interface"
	"backend/internal/application/service"
	"backend/internal/infra/auth"
	"backend/internal/infra/cryptox/hashing"
	"backend/internal/infra/repository"
	service2 "backend/internal/infra/service"
	"context"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"time"
)

type Application struct {
	Config         internal.Config
	PostgresClient *bun.DB

	JwtProvider       *auth.JwtProvider
	ArgonHashProvider *hashing.Argon2idHash

	ProductService _interface.ProductService
	UserService    _interface.UserService
	CacheService   _interface.CacheService
}

func Init(cfg internal.Config) *Application {
	postgresClient, err := InitPostgresConnection(cfg)
	if err != nil {
		log.Panic().Stack().Msgf("fail to connect postgres %+v", err)
	}
	productRepository := repository.NewProductRepositoryImpl(postgresClient)
	cacheService := service2.NewCacheServiceImpl(SetupRedis(cfg))
	productService := service.NewProductServiceImpl(productRepository, cacheService)
	jwtProvider, err := auth.NewJwtProvider(cfg.Auth.JwtPublicKey, cfg.Auth.JwtPrivateKey)
	if err != nil {
		log.Panic().Stack().Msgf("fail to init jwt provider: %+v", err)
	}
	userService := service.NewUserServiceImpl(jwtProvider)

	return &Application{
		Config:         cfg,
		PostgresClient: postgresClient,
		ProductService: productService,
		JwtProvider:    jwtProvider,
		UserService:    userService,
		CacheService:   cacheService,
	}
}

func InitPostgresConnection(appConfig internal.Config) (*bun.DB, error) {
	host := appConfig.Postgres.Host
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", appConfig.Postgres.User, appConfig.Postgres.Password, host, appConfig.Postgres.Database)
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())

	//db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	db.SetMaxIdleConns(appConfig.Postgres.PoolSize)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	err := db.PingContext(ctx)
	return db, err
}

func SetupRedis(cfg internal.Config) (r *redis.Client) {
	var err error
	if !cfg.Redis.UseSentinel {
		opts := &redis.Options{
			Addr:     cfg.Redis.Host,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.Database,
		}
		r, err = setupRedisStandalone(opts)
	} else {
		opts := &redis.FailoverOptions{
			MasterName:       cfg.Redis.MasterName,
			SentinelAddrs:    []string{cfg.Redis.Host},
			SentinelPassword: cfg.Redis.SentinelPassword,

			Password: cfg.Redis.Password,
			DB:       cfg.Redis.Database,
		}
		r, err = setupRedisWithSentinel(opts)
	}
	if err != nil {
		log.Panic().Err(err).Any("cfg", cfg.Redis).Msg("cannot setup redis")
	}
	return
}

func setupRedisWithSentinel(opts *redis.FailoverOptions) (*redis.Client, error) {
	client := redis.NewFailoverClient(opts)
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, errors.Wrap(err, "connect to redis")
	}
	log.Info().Msg("connected to redis")
	return client, nil
}

func setupRedisStandalone(opts *redis.Options) (*redis.Client, error) {
	client := redis.NewClient(opts)
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, errors.Wrap(err, "connect to redis")
	}

	return client, nil
}
