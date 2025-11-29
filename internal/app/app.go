package app

import (
	"net/http"
	"online-subscription/internal/repository"
	"os"
	"time"

	"online-subscription/internal/config"
	"online-subscription/internal/handler"
	"online-subscription/internal/logger"
	"online-subscription/internal/repository/postgres"
	"online-subscription/internal/usecase"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type App struct {
	Server *http.Server
}

func Start() *App {
	cfg := config.LoadConfig(".env")

	if err := logger.Init(); err != nil {
		panic(err)
	}
	defer logger.Sync()

	db, err := repository.ConnectWithRetry(cfg.DSN(), logger.Get(), 10, 2*time.Second)
	if err != nil {
		logger.Error("Failed to connect to DB after retries", zap.Error(err))
		os.Exit(1)
	}

	if err := repository.RunMigrations(db, "file:///app/migrations"); err != nil {
		logger.Error("Failed to run migrations", zap.Error(err))
		os.Exit(1)
	}

	repo := postgres.NewSubscriptionRepo(db)
	uc := usecase.NewSubscriptionUseCase(repo)
	h := handler.NewSubscriptionHandler(uc)

	router := NewRouter(h)

	srv := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: router,
	}

	logger.Info("Starting server", zap.String("port", cfg.AppPort))

	return &App{Server: srv}
}
