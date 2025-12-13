package repository

import (
	"fmt"
	"online-subscription/internal/logger"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func RunMigrations(db *sqlx.DB, migrationsPath string) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		logger.Error("failed to create migration driver", zap.Error(err))
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"postgres", driver)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func ConnectWithRetry(dsn string, logger *zap.Logger, retries int, delay time.Duration) (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error

	for i := 0; i < retries; i++ {
		db, err = sqlx.Open("postgres", dsn)
		if err == nil {
			err = db.Ping()
			if err == nil {
				logger.Info("Connected to database", zap.Int("attempt", i+1))
				return db, nil
			}
		}

		logger.Warn("Database not ready, retrying...", zap.Int("attempt", i+1), zap.Error(err))
		time.Sleep(delay)
	}

	return nil, err
}
