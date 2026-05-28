package main

import (
	"context"
	"net/http"
	_ "online-subscription/docs"
	"online-subscription/internal/app"
	"online-subscription/internal/config"
	"online-subscription/internal/logger"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

// @title Online Subscriptions API service
// @version 1.0
// @description Aggregation data of online subscriptions
// @BasePath /
func main() {
	cfg, err := config.Load(".env")
	if err != nil {
		panic(err)
	}

	log, err := logger.New(cfg.LogLevel)
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	application, err := app.New(cfg, log)
	if err != nil {
		log.Fatal("Failed to initialize application", zap.Error(err))
	}

	srv := application.Server
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server crashed", zap.Error(err))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Info("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Graceful shutdown failed", zap.Error(err))
	} else {
		log.Info("Server stopped gracefully")
	}
}
