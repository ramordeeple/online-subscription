package app

import (
	"fmt"
	"net/http"
	"online-subscription/internal/config"
	"online-subscription/internal/handler"
	"online-subscription/internal/repository"
	"online-subscription/internal/repository/postgres"
	"online-subscription/internal/usecase"
	"time"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type App struct {
	Server *http.Server
}

func New(cfg *config.Config, log *zap.Logger) (*App, error) {
	db, err := repository.ConnectWithRetry(cfg.DSN(), log, 5, 2*time.Second)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	if err := repository.RunMigrations(db, cfg.MigrationsPath); err != nil {
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	repo := postgres.NewSubscriptionRepo(db)
	uc := usecase.NewSubscriptionUseCase(repo)
	h := handler.NewSubscriptionHandler(uc, log)

	srv := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: NewRouter(h),
	}

	log.Info("Starting service", zap.String("port", cfg.AppPort))
	return &App{Server: srv}, nil
}
