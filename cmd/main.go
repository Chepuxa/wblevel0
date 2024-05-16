package main

import (
	"sync"
	"time"
	"wbtech/level0/internal/api"
	"wbtech/level0/internal/api/cache"
	"wbtech/level0/internal/config"
	"wbtech/level0/internal/db"
	"wbtech/level0/internal/logger"
	"wbtech/level0/internal/streaming"
	"wbtech/level0/internal/streaming/subscribers"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"

	"go.uber.org/zap"
)

func main() {
	logger.InitLogger()
	defer logger.CloseLogger()

	envErr := godotenv.Load()
	if envErr != nil {
		logger.Logger.Fatal("Error loading .env file", zap.Error(envErr))
	}

	cfg := config.NewConfig()

	cfgErr := cfg.ParseFlags()
	if cfgErr != nil {
		logger.Logger.Fatal("Failed to parse command-line flags", zap.Error(cfgErr))
	}

	conn, dbErr := db.Connect(cfg)
	if dbErr != nil {
		logger.Logger.Fatal("Failed to connect to the database", zap.Error(dbErr))
		panic(dbErr)
	}
	defer conn.Close()

	migrationErr := db.CreateTables(conn, cfg)
	if migrationErr != nil {
		logger.Logger.Fatal("Failed to apply migrations", zap.Error(migrationErr))
		panic(migrationErr)
	}

	repositories := cfg.InitializeRepositories(conn)

	cache := cache.New(3*time.Hour, 1*time.Hour)
	cacheErr := cache.Fill(*repositories)
	if cacheErr != nil {
		logger.Logger.Fatal("Failed to fill cache", zap.Error(cacheErr))
		panic(cacheErr)
	}

	wg := sync.WaitGroup{}

	va := validator.New()

	handlers := cfg.InitializeHandlers(repositories, cache, &wg, va)

	sc, streamErr := streaming.Connect(cfg)
	if streamErr != nil {
		logger.Logger.Fatal("Failed to connect to nats-streaming server", zap.Error(streamErr))
		panic(streamErr)
	}
	defer sc.Close()

	subs := cfg.InitializeSubscribers(sc, repositories, va, cache, &wg, logger.Logger)
	subscriptions := subscribers.SubscribeAll(subs)
	defer subscribers.UnsubscribeAll(subscriptions)

	srv := api.NewAPI(logger.Logger, cfg, handlers, &wg)
	srvErr := srv.Run()
	if srvErr != nil {
		logger.Logger.Fatal("Failed to start the server", zap.Error(srvErr))
	}
}
