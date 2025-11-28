package app

import (
	"database/sql"
	"net/http"
	"os"

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
	cfg := config.LoadConfig("config.yaml")

	if err := logger.Init(); err != nil {
		panic(err)
	}
	defer logger.Sync()

	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		logger.Error("Failed to connect to DB", zap.Error(err))
		os.Exit(1)
	}

	repo := postgres.NewSubscriptionRepo(db)
	uc := usecase.NewSubscriptionUseCase(repo)
	h := handler.NewSubscriptionHandler(uc)
	router := NewRouter(h)

	srv := &http.Server{
		Addr:    ":" + cfg.App.Port,
		Handler: router,
	}

	logger.Info("Starting server", zap.String("port", cfg.App.Port))

	return &App{Server: srv}
}
