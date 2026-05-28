package repository

import (
	"fmt"
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
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %w", err)
	}

	return nil
}

func ConnectWithRetry(dsn string, log *zap.Logger, retries int, delay time.Duration) (*sqlx.DB, error) {
	var (
		db  *sqlx.DB
		err error
	)

	for i := range retries {
		db, err = sqlx.Open("postgres", dsn)
		if err == nil {
			if err = db.Ping(); err == nil {
				log.Info("Connected to database", zap.Int("attempt", i+1))
				return db, nil
			}
		}
		log.Warn("Database not ready, retrying...", zap.Int("attempt", i+1), zap.Error(err))
		time.Sleep(delay)
	}

	return nil, err
}
