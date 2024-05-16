package db

import (
	"database/sql"
	"errors"
	"log"
	"wbtech/level0/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Connect(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.DSN)

	return db, err
}

func CreateTables(db *sql.DB, cfg *config.Config) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})

	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://internal/db/migrations", cfg.DbName, driver)

	if err != nil {
		return err
	}

	err = m.Up()

	switch {
	case errors.Is(err, migrate.ErrNoChange):
	case err != nil:
		return err
	default:
		log.Println("Migrations applied successfully")
	}
	return nil
}
