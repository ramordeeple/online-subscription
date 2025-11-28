package app

import (
	"database/sql"
	"log"
	"net/http"

	"online-subscription/internal/config"
	"online-subscription/internal/handler"
	"online-subscription/internal/logger"
	"online-subscription/internal/repository/postgres"
	"online-subscription/internal/usecase"

	_ "github.com/lib/pq"
)

type App struct {
	Server *http.Server
}

func Start() *App {
	cfg := config.LoadConfig("config.yaml")

	logg := logger.New()

	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		log.Fatal("failed to connect to db:", err)
	}

	repo := postgres.NewSubscriptionRepo(db)

	uc := usecase.NewSubscriptionUseCase(repo)

	h := handler.NewSubscriptionHandler(uc)

	router := NewRouter(h)

	srv := &http.Server{
		Addr:    ":" + cfg.App.Port,
		Handler: router,
	}

	logg.Info("Starting server on port " + cfg.App.Port)

	return &App{Server: srv}
}
