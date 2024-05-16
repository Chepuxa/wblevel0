package config

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"sync"
	"wbtech/level0/internal/api/cache"
	"wbtech/level0/internal/api/handlers"
	"wbtech/level0/internal/db/repositories"
	"wbtech/level0/internal/streaming/subscribers"

	"github.com/go-playground/validator/v10"
	"github.com/nats-io/stan.go"
	"go.uber.org/zap"
)

type Config struct {
	Address     string
	JwtSecret   string
	DSN         string
	DbName      string
	StanCluster string
	StanClient  string
	StanSubject string
}

func NewConfig() *Config {
	return &Config{}
}

func (cfg *Config) ParseFlags() error {
	address := fmt.Sprintf("%v:%v", os.Getenv("APP_HOST"), os.Getenv("APP_INTERNAL_PORT"))
	flag.StringVar(&cfg.Address, "address", address, "API server address")

	connectionString := fmt.Sprintf("user=%v password=%v host=%v port=%v dbname=%v sslmode=%v",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("DB_INTERNAL_PORT"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_SSL_MODE"))
	flag.StringVar(&cfg.DSN, "DSN", connectionString, "DSN")
	flag.StringVar(&cfg.DbName, "dbName", os.Getenv("POSTGRES_DB"), "DB name")
	flag.StringVar(&cfg.StanCluster, "stanCluster", os.Getenv("STAN_CLUSTER_ID"), "Stan cluster id")
	flag.StringVar(&cfg.StanClient, "stanClient", os.Getenv("STAN_CLIENT_ID"), "Stan client id")
	flag.StringVar(&cfg.StanSubject, "stanSubject", os.Getenv("STAN_SUBJECT"), "Stan subject")
	return nil
}

func (c *Config) InitializeHandlers(r *repositories.Repositories, ca *cache.Cache, wg *sync.WaitGroup, va *validator.Validate) *handlers.Handlers {
	return handlers.NewHandlers(r.OrderRepository, ca, wg, va)
}

func (c *Config) InitializeRepositories(db *sql.DB) *repositories.Repositories {
	return repositories.NewRepositories(db)
}

func (c *Config) InitializeSubscribers(sc stan.Conn, r *repositories.Repositories, va *validator.Validate, ca *cache.Cache,
	wg *sync.WaitGroup, lg *zap.SugaredLogger) *subscribers.Subscribers {
	return subscribers.NewSubscribers(sc, c.StanSubject, r.OrderRepository, va, ca, wg, lg)
}
